// Node 18+, no deps. Run: `node test_api.js [iterations] [concurrency]`
// Exercises full CRUD for /user, /admin and /trainer:
//   signup -> login -> get -> update -> delete
// Also runs auth negative checks (no cookie, garbage cookie, cross-role) expecting 401.
// Reports per-endpoint latency stats and auth-check pass/fail counts.

const BASE = 'http://localhost';
const ITERATIONS = Number(process.argv[2]) || 20;
const CONCURRENCY = Number(process.argv[3]) || 1;

const RESOURCES = ['user', 'admin', 'trainer'];
const STEPS = ['signup', 'login', 'get', 'update', 'delete'];

const samples = {};
const failures = {};
for (const r of RESOURCES) {
  for (const s of STEPS) {
    samples[`${r}_${s}`] = [];
    failures[`${r}_${s}`] = 0;
  }
}
const PUBLIC_STEPS = ['public_subscriptions', 'public_classes'];
for (const s of PUBLIC_STEPS) {
  samples[s] = [];
  failures[s] = 0;
}

async function timed(name, fn) {
  const t0 = performance.now();
  let ok = false;
  let status = 0;
  let bodyText = '';
  try {
    const res = await fn();
    status = res.status;
    bodyText = await res.text();
    ok = res.ok;
    return { res, status, bodyText, ok };
  } finally {
    const ms = performance.now() - t0;
    samples[name].push(ms);
    if (!ok) failures[name]++;
  }
}

function extractCookie(setCookieHeader) {
  if (!setCookieHeader) return null;
  return setCookieHeader.split(';')[0]; // "key=<jwt>"
}

function stats(arr) {
  if (arr.length === 0) return null;
  const sorted = [...arr].sort((a, b) => a - b);
  const pct = (p) => sorted[Math.min(sorted.length - 1, Math.floor(p * sorted.length))];
  const mean = arr.reduce((s, v) => s + v, 0) / arr.length;
  return {
    n: arr.length,
    min: sorted[0],
    mean,
    p50: pct(0.50),
    p95: pct(0.95),
    p99: pct(0.99),
    max: sorted[sorted.length - 1],
  };
}

function fmt(n) { return n == null ? '-' : n.toFixed(2).padStart(8); }

async function setupAccount(resource, i) {
  const email = `perf+${resource}_${Date.now()}_${i}@test.local`;
  const password = 'pw12345';

  const signup = await timed(`${resource}_signup`, () =>
    fetch(`${BASE}/${resource}/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: `${resource}${i}`, email, password }),
    }));
  if (!signup.ok) return null;

  const login = await timed(`${resource}_login`, () =>
    fetch(`${BASE}/${resource}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    }));
  if (!login.ok) return null;

  return extractCookie(login.res.headers.get('set-cookie'));
}

async function runHot(resource, cookie, i) {
  await timed(`${resource}_get`, () =>
    fetch(`${BASE}/${resource}/get`, {
      method: 'GET',
      headers: { cookie },
    }));

  await timed(`${resource}_update`, () =>
    fetch(`${BASE}/${resource}/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', cookie },
      body: JSON.stringify({ name: `${resource}-${i}` }),
    }));
}

async function runPublic() {
  await timed('public_subscriptions', () => fetch(`${BASE}/subscriptions`));
  await timed('public_classes',       () => fetch(`${BASE}/classes`));
}

async function teardownAccount(resource, cookie) {
  await timed(`${resource}_delete`, () =>
    fetch(`${BASE}/${resource}/delete`, {
      method: 'DELETE',
      headers: { cookie },
    }));
}

const PROTECTED = [
  { step: 'get',    method: 'GET',    path: 'get',    body: null },
  { step: 'update', method: 'PUT',    path: 'update', body: { name: 'nope' } },
  { step: 'delete', method: 'DELETE', path: 'delete', body: null },
];

const authChecks = {};
for (const r of RESOURCES) {
  for (const p of PROTECTED) {
    authChecks[`${r}_${p.step}_nocookie`]  = { pass: 0, fail: 0 };
    authChecks[`${r}_${p.step}_badcookie`] = { pass: 0, fail: 0 };
  }
}

async function expectUnauthorized(key, fn) {
  try {
    const res = await fn();
    if (res.status === 401) authChecks[key].pass++;
    else authChecks[key].fail++;
  } catch {
    authChecks[key].fail++;
  }
}

async function runAuthChecks(resource) {
  for (const p of PROTECTED) {
    const url = `${BASE}/${resource}/${p.path}`;
    const init = {
      method: p.method,
      headers: p.body ? { 'Content-Type': 'application/json' } : {},
    };
    if (p.body) init.body = JSON.stringify(p.body);

    await expectUnauthorized(`${resource}_${p.step}_nocookie`, () => fetch(url, init));

    await expectUnauthorized(`${resource}_${p.step}_badcookie`, () => fetch(url, {
      ...init,
      headers: { ...init.headers, cookie: 'key=not-a-real-jwt' },
    }));
  }
}

const crossRoleChecks = {};
for (const from of RESOURCES) {
  for (const to of RESOURCES) {
    if (from === to) continue;
    for (const p of PROTECTED) {
      crossRoleChecks[`${from}_as_${to}_${p.step}`] = { pass: 0, fail: 0 };
    }
  }
}

async function signupAndLogin(resource, tag) {
  const email = `xrole+${resource}_${Date.now()}_${tag}@test.local`;
  const password = 'pw12345';

  const signup = await fetch(`${BASE}/${resource}/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: `${resource}-xrole`, email, password }),
  });
  if (!signup.ok) return null;

  const login = await fetch(`${BASE}/${resource}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  if (!login.ok) return null;

  return extractCookie(login.headers.get('set-cookie'));
}

async function runCrossRoleChecks() {
  const cookies = {};
  for (const r of RESOURCES) cookies[r] = await signupAndLogin(r, 'x');

  for (const from of RESOURCES) {
    const cookie = cookies[from];
    if (!cookie) continue;
    for (const to of RESOURCES) {
      if (from === to) continue;
      for (const p of PROTECTED) {
        const key = `${from}_as_${to}_${p.step}`;
        const init = {
          method: p.method,
          headers: { cookie, ...(p.body ? { 'Content-Type': 'application/json' } : {}) },
        };
        if (p.body) init.body = JSON.stringify(p.body);
        try {
          const res = await fetch(`${BASE}/${to}/${p.path}`, init);
          if (res.status === 401) crossRoleChecks[key].pass++;
          else crossRoleChecks[key].fail++;
        } catch {
          crossRoleChecks[key].fail++;
        }
      }
    }
  }
}

const subChecks = {};
subChecks['list_plans_public']         = { pass: 0, fail: 0 };
subChecks['user_subscribe']            = { pass: 0, fail: 0 };
subChecks['user_subscribe_bad_plan']   = { pass: 0, fail: 0 };
subChecks['user_list_mine']            = { pass: 0, fail: 0 };
subChecks['user_subscribe_nocookie']   = { pass: 0, fail: 0 };
subChecks['user_list_mine_nocookie']   = { pass: 0, fail: 0 };
subChecks['admin_as_user_subscribe']   = { pass: 0, fail: 0 };
subChecks['trainer_as_user_subscribe'] = { pass: 0, fail: 0 };
subChecks['admin_as_user_list_mine']   = { pass: 0, fail: 0 };
subChecks['trainer_as_user_list_mine'] = { pass: 0, fail: 0 };

function record(key, ok) {
  if (ok) subChecks[key].pass++;
  else    subChecks[key].fail++;
}

async function runSubscriptionChecks() {
  const cookies = {};
  for (const r of RESOURCES) cookies[r] = await signupAndLogin(r, 'sub');

  // Plan listing is public.
  {
    const res = await fetch(`${BASE}/subscriptions`);
    record('list_plans_public', res.ok);
  }

  if (cookies.user) {
    // Subscribe to plan 1 (seeded Basic). Body is a raw JSON number to match BindJSON(&uint).
    const ok = await fetch(`${BASE}/user/subscribe`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie: cookies.user },
      body: JSON.stringify(1),
    });
    record('user_subscribe', ok.ok);

    // Bad plan id must not succeed.
    const bad = await fetch(`${BASE}/user/subscribe`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie: cookies.user },
      body: JSON.stringify(9999),
    });
    record('user_subscribe_bad_plan', !bad.ok);

    // Listing my subscriptions must succeed.
    const mine = await fetch(`${BASE}/user/subscription`, { headers: { cookie: cookies.user } });
    record('user_list_mine', mine.ok);
  }

  // No-cookie on user-only subscription endpoints must be 401.
  {
    const res = await fetch(`${BASE}/user/subscribe`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(1),
    });
    record('user_subscribe_nocookie', res.status === 401);
  }
  {
    const res = await fetch(`${BASE}/user/subscription`);
    record('user_list_mine_nocookie', res.status === 401);
  }

  // Cross-role: admin/trainer cookies must not reach user-only endpoints.
  for (const from of ['admin', 'trainer']) {
    const cookie = cookies[from];
    if (!cookie) continue;
    const sub = await fetch(`${BASE}/user/subscribe`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie },
      body: JSON.stringify(1),
    });
    record(`${from}_as_user_subscribe`, sub.status === 401);

    const list = await fetch(`${BASE}/user/subscription`, { headers: { cookie } });
    record(`${from}_as_user_list_mine`, list.status === 401);
  }

  // Cleanup the fresh accounts.
  for (const r of RESOURCES) {
    const cookie = cookies[r];
    if (!cookie) continue;
    await fetch(`${BASE}/${r}/delete`, { method: 'DELETE', headers: { cookie } });
  }
}

const classChecks = {};
classChecks['list_classes_public']           = { pass: 0, fail: 0 };
classChecks['trainer_create_class']          = { pass: 0, fail: 0 };
classChecks['create_class_nocookie']         = { pass: 0, fail: 0 };
classChecks['create_class_bad_body']         = { pass: 0, fail: 0 };
classChecks['user_as_trainer_create_class']  = { pass: 0, fail: 0 };
classChecks['admin_as_trainer_create_class'] = { pass: 0, fail: 0 };

function recordClass(key, ok) {
  if (ok) classChecks[key].pass++;
  else    classChecks[key].fail++;
}

async function runClassChecks() {
  const cookies = {};
  for (const r of RESOURCES) cookies[r] = await signupAndLogin(r, 'class');

  // GET /classes is public (no middleware in router).
  {
    const res = await fetch(`${BASE}/classes`);
    recordClass('list_classes_public', res.ok);
  }

  // Trainer creates a class.
  if (cookies.trainer) {
    const res = await fetch(`${BASE}/trainer/class`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie: cookies.trainer },
      body: JSON.stringify({
        name: 'Yoga test',
        description: 'automated test class',
        date: new Date(Date.now() + 86400000).toISOString(),
        capacity: 10,
      }),
    });
    recordClass('trainer_create_class', res.ok);

    // Missing required fields (no name, no capacity) -> 400 expected.
    const bad = await fetch(`${BASE}/trainer/class`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie: cookies.trainer },
      body: JSON.stringify({ description: 'nothing else' }),
    });
    recordClass('create_class_bad_body', !bad.ok);
  }

  // No cookie -> 401.
  {
    const res = await fetch(`${BASE}/trainer/class`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: 'x', capacity: 1, date: new Date().toISOString() }),
    });
    recordClass('create_class_nocookie', res.status === 401);
  }

  // Cross-role: user/admin cookies must not reach /trainer/class.
  for (const from of ['user', 'admin']) {
    const cookie = cookies[from];
    if (!cookie) continue;
    const res = await fetch(`${BASE}/trainer/class`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', cookie },
      body: JSON.stringify({ name: 'x', capacity: 1, date: new Date().toISOString() }),
    });
    recordClass(`${from}_as_trainer_create_class`, res.status === 401);
  }

  for (const r of RESOURCES) {
    const cookie = cookies[r];
    if (!cookie) continue;
    await fetch(`${BASE}/${r}/delete`, { method: 'DELETE', headers: { cookie } });
  }
}

async function batched(n, concurrency, work) {
  for (let start = 0; start < n; start += concurrency) {
    const batch = [];
    for (let k = 0; k < concurrency && start + k < n; k++) batch.push(work(start + k));
    await Promise.all(batch);
  }
}

(async () => {
  console.log(`Target: ${BASE}`);
  console.log(`Iterations: ${ITERATIONS}  |  Concurrency: ${CONCURRENCY}`);
  console.log(`Resources: ${RESOURCES.join(', ')}\n`);

  for (const resource of RESOURCES) await runAuthChecks(resource);
  await runCrossRoleChecks();
  await runSubscriptionChecks();
  await runClassChecks();

  // Setup: build a pool of `CONCURRENCY` accounts per resource (bcrypt-heavy, not timed as throughput).
  const pools = {};
  process.stdout.write('setup... ');
  const setupStart = performance.now();
  for (const resource of RESOURCES) {
    pools[resource] = new Array(CONCURRENCY);
    await batched(CONCURRENCY, CONCURRENCY, async (i) => {
      pools[resource][i] = await setupAccount(resource, i);
    });
  }
  const setupMs = performance.now() - setupStart;
  console.log(`done (${setupMs.toFixed(0)}ms)`);

  // Hot loop: each iteration does get + update on every resource, reusing pooled cookies (no bcrypt).
  let done = 0;
  const wallStart = performance.now();
  await batched(ITERATIONS, CONCURRENCY, async (i) => {
    for (const resource of RESOURCES) {
      const cookie = pools[resource][i % CONCURRENCY];
      if (cookie) await runHot(resource, cookie, i);
    }
    await runPublic();
    done++;
    process.stdout.write(`\r${done}/${ITERATIONS}`);
  });
  const wallMs = performance.now() - wallStart;
  console.log('\n');

  // Teardown: delete pooled accounts.
  for (const resource of RESOURCES) {
    await batched(pools[resource].length, CONCURRENCY, async (i) => {
      const cookie = pools[resource][i];
      if (cookie) await teardownAccount(resource, cookie);
    });
  }

  const header = 'endpoint        n      min     mean      p50      p95      p99      max    fails';
  console.log(header);
  console.log('-'.repeat(header.length));
  for (const resource of RESOURCES) {
    for (const step of STEPS) {
      const name = `${resource}_${step}`;
      const s = stats(samples[name]);
      if (!s) { console.log(`${name.padEnd(13)}     (no samples)`); continue; }
      console.log(
        `${name.padEnd(13)} ${String(s.n).padStart(4)} ${fmt(s.min)} ${fmt(s.mean)} ${fmt(s.p50)} ${fmt(s.p95)} ${fmt(s.p99)} ${fmt(s.max)} ${String(failures[name]).padStart(8)}`
      );
    }
    console.log('-'.repeat(header.length));
  }
  for (const name of PUBLIC_STEPS) {
    const s = stats(samples[name]);
    if (!s) { console.log(`${name.padEnd(21)} (no samples)`); continue; }
    console.log(
      `${name.padEnd(13)} ${String(s.n).padStart(4)} ${fmt(s.min)} ${fmt(s.mean)} ${fmt(s.p50)} ${fmt(s.p95)} ${fmt(s.p99)} ${fmt(s.max)} ${String(failures[name]).padStart(8)}`
    );
  }
  console.log('-'.repeat(header.length));
  const hotReqs = samples['user_get'].length + samples['user_update'].length
                + samples['admin_get'].length + samples['admin_update'].length
                + samples['trainer_get'].length + samples['trainer_update'].length
                + samples['public_subscriptions'].length + samples['public_classes'].length;
  const totalReqs = Object.values(samples).reduce((s, a) => s + a.length, 0);
  console.log(`\nSetup: ${setupMs.toFixed(0)}ms (${CONCURRENCY * RESOURCES.length} accounts)`);
  console.log(`Hot loop: ${wallMs.toFixed(0)}ms  |  Hot requests: ${hotReqs}  |  Throughput: ${(hotReqs / (wallMs / 1000)).toFixed(1)} req/s`);
  console.log(`Total requests (incl setup/teardown): ${totalReqs}`);

  console.log('\nAuth negative checks (expecting 401):');
  const authHeader = 'check                         pass    fail';
  console.log(authHeader);
  console.log('-'.repeat(authHeader.length));
  for (const resource of RESOURCES) {
    for (const p of PROTECTED) {
      for (const variant of ['nocookie', 'badcookie']) {
        const key = `${resource}_${p.step}_${variant}`;
        const r = authChecks[key];
        console.log(`${key.padEnd(28)} ${String(r.pass).padStart(6)} ${String(r.fail).padStart(7)}`);
      }
    }
  }
  const authTotal = Object.values(authChecks).reduce((s, r) => s + r.pass + r.fail, 0);
  const authFails = Object.values(authChecks).reduce((s, r) => s + r.fail, 0);
  console.log(`\nAuth checks: ${authTotal - authFails}/${authTotal} returned 401 as expected.`);

  console.log('\nCross-role checks (expecting 401):');
  const xHeader = 'check                              pass    fail';
  console.log(xHeader);
  console.log('-'.repeat(xHeader.length));
  for (const from of RESOURCES) {
    for (const to of RESOURCES) {
      if (from === to) continue;
      for (const p of PROTECTED) {
        const key = `${from}_as_${to}_${p.step}`;
        const r = crossRoleChecks[key];
        console.log(`${key.padEnd(33)} ${String(r.pass).padStart(6)} ${String(r.fail).padStart(7)}`);
      }
    }
  }
  const xTotal = Object.values(crossRoleChecks).reduce((s, r) => s + r.pass + r.fail, 0);
  const xFails = Object.values(crossRoleChecks).reduce((s, r) => s + r.fail, 0);
  console.log(`\nCross-role checks: ${xTotal - xFails}/${xTotal} returned 401 as expected.`);

  console.log('\nSubscription checks:');
  const subHeader = 'check                              pass    fail';
  console.log(subHeader);
  console.log('-'.repeat(subHeader.length));
  for (const key of Object.keys(subChecks)) {
    const r = subChecks[key];
    console.log(`${key.padEnd(33)} ${String(r.pass).padStart(6)} ${String(r.fail).padStart(7)}`);
  }
  const subTotal = Object.values(subChecks).reduce((s, r) => s + r.pass + r.fail, 0);
  const subFails = Object.values(subChecks).reduce((s, r) => s + r.fail, 0);
  console.log(`\nSubscription checks: ${subTotal - subFails}/${subTotal} passed.`);

  console.log('\nClass checks:');
  const classHeader = 'check                              pass    fail';
  console.log(classHeader);
  console.log('-'.repeat(classHeader.length));
  for (const key of Object.keys(classChecks)) {
    const r = classChecks[key];
    console.log(`${key.padEnd(33)} ${String(r.pass).padStart(6)} ${String(r.fail).padStart(7)}`);
  }
  const classTotal = Object.values(classChecks).reduce((s, r) => s + r.pass + r.fail, 0);
  const classFails = Object.values(classChecks).reduce((s, r) => s + r.fail, 0);
  console.log(`\nClass checks: ${classTotal - classFails}/${classTotal} passed.`);
})();


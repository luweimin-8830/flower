import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import vm from 'node:vm';

const source = readFileSync(new URL('../subpackages/analytics/charts/charts.vue', import.meta.url), 'utf8')
    .split('<script>')[1].split('</script>')[0]
    .replace(/^import .*$/gm, '').replace('export default', 'globalThis.component =');
function page(request) {
    const context = { callContainer: request, navBar: {}, qiunDataCharts: {}, uni: { showToast() {} }, console: { error() {} } };
    vm.runInNewContext(source, context);
    const instance = context.component.data();
    Object.assign(instance, context.component.methods, { familyId: 1, chartIndex: 2, priceStartDate: '2026-11-01', priceEndDate: '2026-11-17' });
    return instance;
}
const stats = { totalCents: 1234, pricedCount: 1, unpricedCount: 2, buckets: [] };
test('price date filter sends selected family and dates', async () => {
    let args;
    const p = page(async (...a) => { args = a; return { code: 0, data: stats }; });
    await p.handlePriceDateChange('start', { detail: { value: '2026-10-01' } });
    assert.equal(args[0], '/api/chart/price');
    assert.deepEqual(JSON.parse(JSON.stringify(args[1])), { familyId: 1, startDate: '2026-10-01', endDate: '2026-11-17' });
    assert.equal(p.priceStats.totalCents, 1234);
    await p.handlePriceDateChange('start', { detail: { value: '2026-12-01' } });
    assert.equal(p.priceStartDate, '2026-10-01');
});
test('older response cannot overwrite newer date range', async () => {
    const pending = [];
    const p = page(() => new Promise(resolve => pending.push(resolve)));
    const old = p.getServerData();
    const latest = p.handlePriceDateChange('end', { detail: { value: '2026-11-18' } });
    pending[1]({ code: 0, data: stats }); await latest;
    pending[0]({ code: 0, data: { ...stats, totalCents: 99999 } }); await old;
    assert.equal(p.priceStats.totalCents, 1234);
    assert.equal(p.priceLoading, false);
});
test('failed API clears old totals and supports retry', async () => {
    let fail = true;
    const p = page(async () => fail ? { code: -1, msg: 'failure' } : { code: 0, data: stats });
    p.priceStats = stats;
    await p.getServerData();
    assert.equal(p.priceStats, null);
    assert.ok(p.priceError);
    fail = false; await p.getServerData();
    assert.equal(p.priceError, '');
    assert.equal(p.priceStats.totalCents, 1234);
});

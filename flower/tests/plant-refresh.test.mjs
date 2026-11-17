import { readFileSync } from 'node:fs';
import assert from 'node:assert/strict';
import test from 'node:test';
const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;
for (const file of ['components/home.vue', 'pages/sharePlantList/sharePlantList.vue']) {
    const source = readFileSync(new URL('../' + file, import.meta.url), 'utf8');
    const body = source.split('async getPlantsList(isLoadMore = false) {')[1].split('onScrollToLower()')[0].replace(/},\s*$/, '');
    const requests = [];
    const callContainer = () => new Promise(resolve => requests.push(resolve));
    const run = new AsyncFunction('callContainer', 'isLoadMore', body);
    const create = () => ({ plantsList: [{ id: 1, cover: { url: 'old' } }], value: 1, familyId: 1,
        tagList: [{ id: 0 }, { id: 2 }], currentTagIndex: 0, page: 1, pageSize: 20, total: 1, isLoading: false });
    const response = url => ({ data: { list: [{ id: 1, cover: { url } }], total: 1 } });
    test(file + ': same first ID still updates the cover', async () => {
        const vm = create();
        const refresh = run.call(vm, callContainer, false);
        requests.shift()(response('new'));
        await refresh;
        assert.equal(vm.plantsList[0].cover.url, 'new');
    });
    test(file + ': refresh during loading wins over older responses', async () => {
        const vm = create();
        const first = run.call(vm, callContainer, false);
        vm.currentTagIndex = 1;
        const second = run.call(vm, callContainer, false);
        const resolveFirst = requests.shift(), resolveSecond = requests.shift();
        assert.equal(typeof resolveSecond, 'function');
        resolveSecond(response('newest')); await second;
        resolveFirst(response('stale')); await first;
        assert.equal(vm.plantsList[0].cover.url, 'newest');
        assert.equal(vm.isLoading, false);
    });
}

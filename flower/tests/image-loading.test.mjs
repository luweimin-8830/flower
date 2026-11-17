import { readFileSync } from 'node:fs';
import assert from 'node:assert/strict';
import test from 'node:test';
const moduleURL = source => 'data:text/javascript;base64,' + Buffer.from(source).toString('base64');
const queueURL = moduleURL(readFileSync(new URL('../utils/imageLoadQueue.js', import.meta.url), 'utf8'));
const thumbnailURL = moduleURL(readFileSync(new URL('../utils/imageThumbnail.js', import.meta.url), 'utf8'));
const { createImageLoadQueue } = await import(queueURL);

test('caps requests, skips canceled work, and releases each slot only once', () => {
    const enqueue = createImageLoadQueue(4), started = [], done = [], cancel = [];
    for (let i = 0; i < 20; i++) cancel.push(enqueue(finish => { started.push(i); done[i] = finish; }));
    assert.deepEqual(started, [0, 1, 2, 3]);
    cancel[4](); done[1]();
    assert.deepEqual(started, [0, 1, 2, 3, 5]);
    done[1](); assert.equal(started.length, 5);
    cancel[0](); assert.equal(started.at(-1), 6);
    cancel.forEach(finish => finish());
    let resumed = false;
    enqueue(finish => { resumed = true; finish(); });
    assert.ok(resumed);
});

test('offscreen images wait; visibility, error, retry, URL changes and disposal work', async () => {
    const source = readFileSync(new URL('../components/LazyPlantImage.vue', import.meta.url), 'utf8')
        .split('<script>')[1].split('</script>')[0]
        .replace("'../utils/imageLoadQueue'", JSON.stringify(queueURL))
        .replace("'../utils/imageThumbnail'", JSON.stringify(thumbnailURL));
    const { default: component } = await import(moduleURL(source));
    let notify, disconnected = false;
    globalThis.uni = { getWindowInfo: () => ({ pixelRatio: 3 }), createIntersectionObserver: () => ({
        relativeToViewport() { return this; },
        observe(selector, callback) { notify = callback; },
        disconnect() { disconnected = true; }
    }) };
    const vm = { ...component.data(), src: 'cloud://test/a.jpg', $nextTick: fn => fn() };
    for (const [key, fn] of Object.entries(component.methods)) vm[key] = fn.bind(vm);
    component.mounted.call(vm);
    assert.equal(vm.source, '');
    notify({ intersectionRatio: 1 });
    await new Promise(resolve => setImmediate(resolve));
    assert.equal(vm.source, vm.src);
    const event = () => ({ currentTarget: { dataset: { request: `${vm.generation}:${vm.source}` } } });
    const staleEvent = event();
    vm.onLoad(event()); assert.equal(vm.status, 'loaded');
    vm.src = 'cloud://test/b.jpg'; component.watch.src.call(vm);
    await new Promise(resolve => setImmediate(resolve));
    assert.equal(vm.source, vm.src);
    vm.onLoad(staleEvent); assert.equal(vm.status, 'loading');
    vm.onError(staleEvent); assert.equal(vm.source, vm.src);
    vm.onError(event()); assert.equal(vm.source, ''); assert.equal(vm.status, 'error');
    vm.retry(); assert.equal(vm.status, 'loading');
    component.beforeUnmount.call(vm);
    assert.ok(disconnected); assert.equal(vm.source, '');
    vm.request(); assert.equal(vm.source, '');
    delete globalThis.uni;
});

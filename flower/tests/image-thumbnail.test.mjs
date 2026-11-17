import { readFileSync } from 'node:fs';
import assert from 'node:assert/strict';
import test from 'node:test';
const source = readFileSync(new URL('../utils/imageThumbnail.js', import.meta.url), 'utf8');
const { thumbnailWidth, thumbnailSource } = await import('data:text/javascript;base64,' + Buffer.from(source).toString('base64'));
test('physical pixels cover the crop without enlarging originals', () => {
    assert.equal(thumbnailWidth({ width: 180, height: 270 }, 3, 1080, 1920), 544);
    assert.equal(thumbnailWidth({ width: 180, height: 270 }, 3, 1920, 1080), 1440);
    assert.equal(thumbnailWidth({ width: 180, height: 270 }, 3, 320, 480), 320);
    assert.equal(thumbnailWidth({ width: 56, height: 56 }, 3, 1080, 1920), 192);
    assert.equal(thumbnailWidth(null, 3), 0);
});
test('resolve cloud URLs once, preserve authorization queries, skip unrelated URLs', async () => {
    let calls = 0;
    const resolve = async () => { calls++; return 'https://example.tcb.qcloud.la/a.jpg?token=abc'; };
    const result = await thumbnailSource('cloud://test/image', 544, resolve);
    assert.equal(result, 'https://example.tcb.qcloud.la/a.jpg?token=abc&imageMogr2/thumbnail/544x/quality/80');
    await thumbnailSource('cloud://test/image', 192, resolve);
    assert.equal(calls, 1);
    for (const url of ['/static/default.svg', 'https://other.example/a.jpg?signature=123']) {
        assert.equal(await thumbnailSource(url, 544, resolve), url);
    }
    const transformed = 'https://example.tcb.qcloud.la/a.jpg?imageMogr2/thumbnail/100x';
    assert.equal(await thumbnailSource(transformed, 544, resolve), transformed);
});

// Request physical pixels for the rendered box, including aspectFill cropping.
export function thumbnailWidth(box, pixelRatio, originalWidth, originalHeight) {
    const width = Number(box?.width), height = Number(box?.height);
    if (!(width > 0) || !(height > 0)) return 0;
    const ratio = Number(pixelRatio) > 0 ? Number(pixelRatio) : 1;
    const ow = Number(originalWidth), oh = Number(originalHeight);
    let target = width * ratio;
    if (ow > 0 && oh > 0) target = Math.max(target, height * ratio * ow / oh);
    // Small buckets improve CDN reuse without undersizing the image.
    target = Math.ceil(target / 32) * 32;
    return Math.max(1, Math.min(10000, ow > 0 ? ow : 10000, target));
}

const links = new Map();
export async function thumbnailSource(src, width, resolveCloudFile) {
    if (!width || !src) return src;
    let url = src;
    if (src.startsWith('cloud://')) {
        let entry = links.get(src);
        if (!entry || entry.expires <= Date.now()) {
            const promise = Promise.resolve().then(() => resolveCloudFile(src));
            entry = { promise, expires: Date.now() + 120000 };
            links.set(src, entry);
            if (links.size > 200) links.delete(links.keys().next().value);
            promise.catch(() => { if (links.get(src) === entry) links.delete(src); });
        }
        url = await entry.promise;
    } else if (!/^https:\/\/[^/?#]+\.(?:tcb\.qcloud\.la|tcloudbaseapp\.com|myqcloud\.com)\//i.test(src)) {
        return src;
    }
    // Preserve existing transformations and never alter local assets or file IDs.
    if (!/^https:\/\//i.test(url) || /[?&](?:imageMogr2|imageView2|ci-process)/i.test(url)) return url;
    const [base, hash] = url.split('#');
    return `${base}${base.includes('?') ? '&' : '?'}imageMogr2/thumbnail/${width}x/quality/80${hash ? '#' + hash : ''}`;
}

// Yuan in the form; integer cents in the API and database. null means unknown.
export function parsePlantPrice(value) {
    const text = String(value ?? '').trim();
    if (!text) return null;
    if (!/^\d{1,6}(\.\d{1,2})?$/.test(text)) throw new Error('价格须为 0–999999.99 元，最多两位小数');
    const [yuan, fraction = ''] = text.split('.');
    return Number(yuan) * 100 + Number(fraction.padEnd(2, '0'));
}
export function formatPlantPrice(cents) {
    if (cents === null || cents === undefined) return '';
    return (cents / 100).toFixed(2);
}

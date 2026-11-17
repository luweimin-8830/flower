// Shared by list images: only visible/near-visible images enter this queue.
export function createImageLoadQueue(limit = 4) {
    let active = 0;
    const pending = [];
    const pump = () => {
        while (active < limit && pending.length) {
            const job = pending.shift();
            if (job.done) continue;
            job.active = true;
            active++;
            job.start(job.finish);
        }
    };
    return (start) => {
        const job = { start, active: false, done: false };
        job.finish = () => {
            if (job.done) return;
            job.done = true;
            if (job.active) active--;
            else {
                const index = pending.indexOf(job);
                if (index >= 0) pending.splice(index, 1);
            }
            pump();
        };
        pending.push(job);
        pump();
        return job.finish;
    };
}
export const enqueueImage = createImageLoadQueue(4);

import { base } from "$app/paths";
import { bidStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export function load(event: LoadEvent) {
    const subscribeToBids = () => {
        const bidInfo = get(bidStore);
        if (bidInfo === undefined) {
            return;
        }
        const sse = new EventSource(
            `${base}/api/bidding/${bidInfo.id}/subscribe`,
        );
        return sse;
    };

    const submitBid = async (amount: number) => {
        const bidInfo = get(bidStore);
        if (bidInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            `${base}/api/bidding/${bidInfo.id}`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ amount })
            },
        );
        return [response.ok, response.status];
    };

    return {
        subscribeToBids,
        submitBid,
    };
}
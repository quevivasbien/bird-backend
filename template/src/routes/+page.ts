import { base } from "$app/paths";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const testAuth = async () => {
        const response = await event.fetch(
            base + "/api/login/testAuth",
        );
        console.log(response);
        const text = await response.text();
        console.log(text);
        return text;
    }
    return {
        testAuth,
    }
}
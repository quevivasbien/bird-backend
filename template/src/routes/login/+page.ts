import { base } from "$app/paths";
import { userStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const login = async (name: string, password: string) => {
        const response = await event.fetch(
            base + "api/auth/login",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name,
                    password,
                }),
            }
        );
        console.log(response);
        if (response.ok) {
            const userInfo = await response.json();
            userStore.set(userInfo);
        }
        return [response.ok, response.status];
    }
    
    return {
        login,
    };
}
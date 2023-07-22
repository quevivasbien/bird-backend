import { base } from "$app/paths";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const register = async (name: string, password: string) => {
        const response = await event.fetch(
            base + "api/auth/register",
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
        return [response.ok, response.status];
    }
    return {
        register,
    };
}
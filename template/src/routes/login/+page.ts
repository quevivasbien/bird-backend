import { base } from "$app/paths";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const login = async (name: string, password: string) => {
        const response = await event.fetch(
            base + "api/login",
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
        return response.status;
    }
    
    return {
        login,
    }    
}
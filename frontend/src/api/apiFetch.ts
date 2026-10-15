export async function apiFetch (
    url: string,
    options: RequestInit = {},
) {
    try {
        const response = await fetch(url, {
            ...options,
            credentials: "include",
            headers: {
                ...options.headers,
            },
        });
        
        const data = await response.json();

        return {
          status: response.status,
          data,
        };
    } catch {
        return null;
    }
}
export async function apiFetch (
    url: string,
    options: RequestInit = {},
    accessToken?: string | null
) {
    try {
        const response = await fetch(url, {
            ...options,
            credentials: "include",
            headers: {
                ...options.headers,
                ...(accessToken
                    ? { Authorization: accessToken }
                    : {}),
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

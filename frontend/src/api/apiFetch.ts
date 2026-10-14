// export async function apiFetch (
//     url: string,
//     options: RequestInit = {},
// ) {
//     try {
//         const response = await fetch(url, {
//             ...options,
//             credentials: "include",
//             headers: {
//                 ...options.headers,
//             },
//         });
        
//         const data = await response.json();

//         return {
//           status: response.status,
//           data,
//         };
//     } catch {
//         return null;
//     }
// }

export async function apiFetch(
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

    const contentType = response.headers.get("content-type");
    const data = contentType?.includes("application/json")
      ? await response.json()
      : await response.text();

    return {
      status: response.status,
      data,
    };
  } catch (error) {
    console.error("API request failed:", error);
    return null;
  }
}

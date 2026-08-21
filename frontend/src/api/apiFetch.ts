export async function apiFetch (
  url: string,
  options: RequestInit
) {
  try {
    const response = await fetch(url, options);
    const data = await response.json();

    return {
      status: response.status,
      data,
    };
  } catch {
    return null;
  }
}
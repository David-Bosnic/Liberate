export interface LiberateRequest {
  prompt: string;
}

export interface LinkRating {
  url: string;
  rating: string;
}

export interface LiberateResponse {
  links: LinkRating[];
}

const url = "http://localhost:8081/api";
const MOCK = false;

export async function CallLiberate(prompt: string): Promise<LiberateResponse> {
  if (MOCK) {
    await new Promise((r) => setTimeout(r, 1000));
    return {
      links: [
        { url: "http://apple.com", rating: "4" },
        { url: "http://corn.com", rating: "3" },
        { url: "http://taco.com", rating: "5" },
      ],
    };
  }
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ prompt: prompt } satisfies LiberateRequest),
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status} ${res.statusText}`);
  }

  return (await res.json()) as LiberateResponse;
}

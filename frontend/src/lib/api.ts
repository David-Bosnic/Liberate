export interface LiberateRequest {
  prompt: string;
}

export interface LiberateResponse {
  links: string[];
}

const url = "http://localhost:8081/api";
const MOCK = false;

export async function CallLiberate(prompt: string): Promise<LiberateResponse> {
  if (MOCK) {
    await new Promise((r) => setTimeout(r, 1000));
    return {
      links: ["http://apple.com", "http://corn.com", "http://taco.com"],
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

export interface LiberateRequest {
  Prompt: string;
}

export interface LiberateResponse {
  Links: string[];
}

const url = "http://localhost:8081/api";
const MOCK = true;

export async function CallLiberate(prompt: string): Promise<LiberateResponse> {
  if (MOCK) {
    await new Promise((r) => setTimeout(r, 1000));
    return {
      Links: ["http://apple.com", "http://corn.com", "http://taco.com"],
    };
  }
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Prompt: prompt } satisfies LiberateRequest),
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status} ${res.statusText}`);
  }

  return res.json() as Promise<LiberateResponse>;
}

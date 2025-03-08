import type { GetMostUsedDomainsRequest, GetMostUsedDomainsResponse } from "~/@types/types";

export const getMostUsedDomains = async (
  request: GetMostUsedDomainsRequest,
): Promise<GetMostUsedDomainsResponse[]> => {
  const resp = await fetch(`/api/v1/dashboard/most-used-domains`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  return await resp.json();
};

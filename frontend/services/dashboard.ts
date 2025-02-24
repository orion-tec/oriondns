export const getMostUsedDomains = async (selectedRange: string, selectedCategories: string[]) => {
  const resp = await fetch(`/api/v1/dashboard/most-used-domains`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      range: selectedRange,
      categories: selectedCategories,
    }),
  });

  return await resp.json();
};

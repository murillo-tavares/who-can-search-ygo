import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { SearchPage } from "./SearchPage";

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <SearchPage />
    </QueryClientProvider>,
  );
}

describe("SearchPage", () => {
  it("renders the initial search experience", () => {
    renderPage();

    expect(screen.getByRole("searchbox", { name: "Target card" })).toBeInTheDocument();
    expect(screen.getByText("Who can")).toBeInTheDocument();
    expect(screen.getByText("add")).toBeInTheDocument();
    expect(screen.queryByLabelText("Searcher results")).not.toBeVisible();
  });
});

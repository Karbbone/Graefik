import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { TasksPanel } from "@/features/tasks";
import { server } from "@/test/msw/server";

describe("TasksPanel", () => {
  it("affiche un état vide quand il n’y a aucune tâche", async () => {
    render(<TasksPanel />);

    expect(
      await screen.findByText(/aucune tâche pour le moment/i),
    ).toBeInTheDocument();
  });

  it("affiche les tâches renvoyées par l’API", async () => {
    server.use(
      http.get("/api/tasks", () =>
        HttpResponse.json([
          { id: "1", title: "Tâche existante", done: false, createdAt: "" },
        ]),
      ),
    );

    render(<TasksPanel />);

    expect(await screen.findByText("Tâche existante")).toBeInTheDocument();
  });

  it("crée une nouvelle tâche via le formulaire", async () => {
    const user = userEvent.setup();
    render(<TasksPanel />);

    // On attend le chargement initial (liste vide).
    await screen.findByText(/aucune tâche pour le moment/i);

    await user.type(
      screen.getByRole("textbox", { name: /titre de la tâche/i }),
      "Ma nouvelle tâche",
    );
    await user.click(screen.getByRole("button", { name: /ajouter/i }));

    expect(await screen.findByText("Ma nouvelle tâche")).toBeInTheDocument();
  });
});

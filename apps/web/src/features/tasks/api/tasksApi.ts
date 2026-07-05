import { apiFetch } from "@/shared/lib/http";
import type { Task } from "../model/task";

export function listTasks(): Promise<Task[]> {
  return apiFetch<Task[]>("/api/tasks");
}

export function createTask(title: string): Promise<Task> {
  return apiFetch<Task>("/api/tasks", {
    method: "POST",
    body: JSON.stringify({ title }),
  });
}

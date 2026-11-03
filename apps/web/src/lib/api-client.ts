import { createClient } from "@/lib/supabase/client";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/v1";

export async function fetchFromApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const supabase = createClient();
  const { data: { session } } = await supabase.auth.getSession();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...options?.headers as Record<string, string>,
  };

  if (session?.access_token) {
    headers["Authorization"] = `Bearer ${session.access_token}`;
  }

  const res = await fetch(`${API_URL}${endpoint}`, { ...options, headers });

  // Leemos como texto para evitar errores si el backend devuelve vacío o texto plano
  const text = await res.text();

  if (!res.ok) {
    let errorMessage = "Error al conectar con el servidor";
    let errorCode = "unknown";
    
    try {
      if (text) {
        const errorData = JSON.parse(text);
        errorMessage = errorData.message || errorMessage;
        errorCode = errorData.error || errorCode;
      }
    } catch {
      errorMessage = text || res.statusText;
    }

    const err = new Error(errorMessage);
    (err as any).code = errorCode;
    throw err;
  }

  if (!text || res.status === 204) {
    return undefined as T;
  }

  try {
    return JSON.parse(text) as T;
  } catch {
    throw new Error("El servidor devolvió un formato inválido.");
  }
}
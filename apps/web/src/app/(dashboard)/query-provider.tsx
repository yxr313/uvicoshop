"use client";

import { QueryClient, QueryClientProvider, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState, useMemo, useRef } from "react";
import { toast } from "sonner";
import { createClient } from "@/lib/supabase/client";

// Componente interno para manejar los efectos secundarios (Realtime, Audio)
function GlobalListeners() {
  const queryClient = useQueryClient();
  const supabase = useMemo(() => createClient(), []);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  useEffect(() => {
    // 1. Desbloqueo de audio con la primera interacción
    audioRef.current = new Audio("/sounds/sale.mp3");
    audioRef.current.volume = 0.5;
    audioRef.current.load();

    const unlock = () => {
      if (!audioRef.current) return;
      audioRef.current.play().then(() => {
        audioRef.current?.pause();
        if (audioRef.current) audioRef.current.currentTime = 0;
      }).catch(() => {});
      document.removeEventListener("click", unlock);
      document.removeEventListener("keydown", unlock);
      document.removeEventListener("touchstart", unlock);
    };

    document.addEventListener("click", unlock);
    document.addEventListener("keydown", unlock);
    document.addEventListener("touchstart", unlock);

    return () => {
      document.removeEventListener("click", unlock);
      document.removeEventListener("keydown", unlock);
      document.removeEventListener("touchstart", unlock);
    };
  }, []);

  useEffect(() => {
    // 2. Suscripción Realtime Global
    // Usamos un nombre estático ya que solo habrá una instancia en toda la app
    const channel = supabase
      .channel("global-realtime-orders")
      .on(
        "postgres_changes",
        { event: "INSERT", schema: "public", table: "orders" },
        (payload) => {
          const newOrder = payload.new as { total_cents: number };
          
          // Reproducir audio
          if (audioRef.current) {
            audioRef.current.play().catch(() => {});
          }

          // Notificación visual
          toast.success("¡Nueva venta recibida!", {
            description: `Pedido por $${(newOrder.total_cents / 100).toFixed(2)}`,
            duration: 8000,
          });

          // Refrescar todos los datos
          queryClient.invalidateQueries({ queryKey: ["dashboard-orders"] });
          queryClient.invalidateQueries({ queryKey: ["dashboard-traffic"] });
          queryClient.invalidateQueries({ queryKey: ["ecommerce-operations"] });
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [supabase, queryClient]);

  return null;
}

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,
        refetchOnWindowFocus: false,
      },
    },
  }));

  return (
    <QueryClientProvider client={queryClient}>
      <GlobalListeners />
      {children}
    </QueryClientProvider>
  );
}
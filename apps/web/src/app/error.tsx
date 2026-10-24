"use client";

import { useEffect } from "react";
import Link from "next/link";
import { ArrowLeft, RefreshCcw } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function ErrorPage({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 bg-background">
      <div className="max-w-md w-full text-center">
        <p className="text-[6rem] leading-none font-display font-medium tracking-tighter text-neutral-200 mb-4">
          500
        </p>
        <h1 className="text-display-3 font-display text-foreground mb-4">
          Algo salió mal.
        </h1>
        <p className="text-body text-muted-foreground mb-10">
          Hubo un error inesperado en nuestro servidor. Nuestro equipo técnico ha sido notificado. Por favor, intenta de nuevo.
        </p>
        
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button onClick={() => reset()} className="w-full sm:w-auto rounded-full bg-accent text-accent-foreground hover:bg-accent-hover">
            <RefreshCcw className="mr-2" size={18} /> Intentar de nuevo
          </Button>
          <Link href="/" className="w-full sm:w-auto">
            <Button variant="outline" className="w-full rounded-full">
              <ArrowLeft className="mr-2" size={18} /> Volver al inicio
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}

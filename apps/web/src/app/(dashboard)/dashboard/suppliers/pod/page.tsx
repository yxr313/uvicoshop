"use client";

import { PageHeader } from "@/components/page-header";
import { Card } from "@/components/ui/card";

export default function PODDesignStudioPage() {
  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader
        title="Estudio de Diseño (POD)"
        description="Sube tu arte, elige las prendas y personaliza tus productos."
      />
      <Card className="p-12 flex flex-col items-center justify-center text-center min-h-[400px] border-dashed">
        <h3 className="text-lg font-medium text-foreground mb-2">Editor de Mockups</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          El editor visual para subir logos, elegir telas y enviar a aprobación del proveedor estará aquí.
        </p>
      </Card>
    </div>
  );
}
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { IconPlus } from "@tabler/icons-react";

export default function StudioHomePage() {
  return (
    <div className="flex flex-col h-full w-full max-w-6xl mx-auto p-8 md:p-12">
      
      {/* Header Editorial */}
      <div className="flex flex-col md:flex-row md:items-end justify-between mb-12 gap-4 border-b border-border pb-8">
        <div>
          <h1 className="text-4xl font-display font-medium text-foreground tracking-tight">
            Mis Diseños
          </h1>
          <p className="text-base text-muted-foreground mt-2 max-w-lg">
            Gestiona y edita tus creaciones para Print on Demand. Sube tus imágenes, añade texto y prepara todo para vender.
          </p>
        </div>
        <Button asChild className="bg-[#95B54D] text-white hover:bg-[#A7C957] h-11 px-6 rounded-lg font-medium">
          <Link href="/editor/new">
            <IconPlus className="size-4 mr-2" /> Nuevo Diseño
          </Link>
        </Button>
      </div>

      {/* Grid de Diseños */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        
        {/* Tarjeta de Estado Vacío (o tu primer diseño) */}
        <Link 
          href="/editor/new"
          className="aspect-[4/5] border border-dashed border-border rounded-2xl flex flex-col items-center justify-center text-center p-8 bg-white hover:border-[#95B54D] hover:bg-[#F0F5E6]/30 transition-all group"
        >
          <div className="size-16 rounded-full bg-muted flex items-center justify-center mb-4 group-hover:bg-[#F0F5E6] transition-colors">
            <IconPlus className="size-8 text-muted-foreground group-hover:text-[#95B54D] transition-colors" />
          </div>
          <h3 className="font-display font-medium text-lg text-foreground">
            Crear mi primer diseño
          </h3>
          <p className="text-sm text-muted-foreground mt-2">
            Empieza a diseñar tu producto ideal.
          </p>
        </Link>

      </div>
    </div>
  );
}
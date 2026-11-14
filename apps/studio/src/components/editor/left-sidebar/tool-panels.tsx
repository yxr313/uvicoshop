"use client";

import { useEditorStore } from "@/store/use-editor-store";
import { Button } from "@/components/ui/button";
import { IconPhoto, IconTypography } from "@tabler/icons-react";

export function ToolPanels() {
    const { activeTool } = useEditorStore();

    if (!activeTool) return null;

    return (
        <div className="w-64 bg-white border-r border-border flex flex-col shrink-0 z-30">
            <div className="h-14 border-b border-border flex items-center px-4">
                <h2 className="font-display font-medium text-foreground capitalize">
                    {activeTool === "ai" ? "Herramientas IA" : activeTool}
                </h2>
            </div>

            <div className="p-4 flex flex-col gap-3">
                {activeTool === "image" && (
                    <Button variant="outline" className="w-full justify-start h-12">
                        <IconPhoto className="size-4 mr-3" /> Subir Imagen
                    </Button>
                )}

                {activeTool === "text" && (
                    <Button variant="outline" className="w-full justify-start h-12">
                        <IconTypography className="size-4 mr-3" /> Añadir Texto
                    </Button>
                )}

                {/* Aquí irán los componentes de Imagen, Texto, etc. */}
                <div className="text-sm text-muted-foreground mt-4 text-center">
                    Opciones de {activeTool} en construcción.
                </div>
            </div>
        </div>
    );
}
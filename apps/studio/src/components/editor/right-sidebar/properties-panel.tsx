"use client";

import { useEditorStore } from "@/store/use-editor-store";

export function RightSidebar() {
    const { selectedId } = useEditorStore();

    return (
        <div className="w-72 bg-white border-l border-border flex flex-col shrink-0 z-40">
            <div className="h-14 border-b border-border flex items-center px-4">
                <h2 className="font-display font-medium">Propiedades</h2>
            </div>
            <div className="p-4 flex-1 overflow-y-auto">
                {selectedId ? (
                    <p className="text-sm text-muted-foreground">
                        Editando elemento seleccionado...
                    </p>
                ) : (
                    <p className="text-sm text-muted-foreground text-center mt-10">
                        Selecciona un elemento en el lienzo para editar sus propiedades.
                    </p>
                )}
            </div>
        </div>
    );
}
"use client";

import { useEditorStore } from "@/store/use-editor-store";

export function LayersPanel() {
    const { elements, selectedId, setSelectedId } = useEditorStore();

    if (elements.length === 0) {
        return (
            <div className="p-6 text-center text-sm text-muted-foreground">
                No hay capas en el diseño.
            </div>
        );
    }

    return (
        <div className="flex flex-col py-2">
            {[...elements].reverse().map((el) => (
                <div
                    key={el.id}
                    onClick={() => setSelectedId(el.id)}
                    className={`px-4 py-3 text-sm cursor-pointer border-l-2 transition-colors ${selectedId === el.id
                            ? "bg-[#F0F5E6] border-[#95B54D] text-foreground"
                            : "border-transparent text-muted-foreground hover:bg-muted/50"
                        }`}
                >
                    {el.type === "text" ? (el.text || "Texto") : "Imagen"}
                </div>
            ))}
        </div>
    );
}
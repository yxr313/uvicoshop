"use client";

import { useEditorStore } from "@/store/use-editor-store";
import { Button } from "@/components/ui/button";
import { IconTypography, IconPhoto, IconWand, IconShape } from "@tabler/icons-react";
import { cn } from "@/lib/utils";

const TOOLS = [
    { id: "text", icon: IconTypography, label: "Texto" },
    { id: "image", icon: IconPhoto, label: "Imágenes" },
    { id: "shapes", icon: IconShape, label: "Formas" },
    { id: "ai", icon: IconWand, label: "Herramientas IA" },
] as const;

export function EditorSidebar() {
    const { activeTool, setActiveTool } = useEditorStore();

    return (
        <div className="w-16 bg-white border-r border-border flex flex-col items-center py-4 gap-2 shrink-0 z-40">
            {TOOLS.map((tool) => (
                <Button
                    key={tool.id}
                    variant="ghost"
                    size="icon"
                    onClick={() => setActiveTool(activeTool === tool.id ? null : tool.id)}
                    className={cn(
                        "size-12 rounded-xl transition-colors",
                        activeTool === tool.id ? "bg-[#F0F5E6] text-[#95B54D] hover:bg-[#F0F5E6]" : "text-muted-foreground"
                    )}
                    title={tool.label}
                >
                    <tool.icon className="size-5" />
                </Button>
            ))}
        </div>
    );
}
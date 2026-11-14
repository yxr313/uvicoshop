"use client";

import { Button } from "@/components/ui/button";
import { IconDeviceFloppy, IconArrowLeft } from "@tabler/icons-react";
import Link from "next/link";

export function EditorHeader({ designId }: { designId: string }) {
    return (
        <header className="h-14 bg-white border-b border-border flex items-center justify-between px-4 shrink-0 z-50">
            <div className="flex items-center gap-4">
                <Button variant="ghost" size="icon" asChild>
                    <Link href="/"><IconArrowLeft className="size-5" /></Link>
                </Button>
                <span className="font-display font-medium">Editor de Diseño</span>
            </div>

            <div className="flex items-center gap-2">
                {/* Botones de Zoom y Deshacer irán aquí */}
                <Button className="bg-[#95B54D] text-white hover:bg-[#A7C957]">
                    <IconDeviceFloppy className="size-4 mr-2" />
                    Guardar
                </Button>
            </div>
        </header>
    );
}
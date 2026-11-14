"use client";

export function CanvasWrapper() {
    return (
        <div className="flex-1 relative flex items-center justify-center overflow-hidden p-6">
            {/* Área de trabajo visual */}
            <div className="w-[500px] h-[600px] bg-white border border-border shadow-sm flex items-center justify-center relative">
                <p className="text-muted-foreground text-sm">
                    [ Aquí se integrará el motor de dibujo ]
                </p>
            </div>
        </div>
    );
}
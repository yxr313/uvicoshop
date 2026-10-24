"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface SlugBuilderProps {
  variant?: "default" | "dark";
  className?: string;
}

export function SlugBuilder({ variant = "dark", className }: SlugBuilderProps) {
  const [slug, setSlug] = useState("");

  return (
    <div 
      className={cn(
        "flex flex-col sm:flex-row items-center w-full p-1.5 bg-white border border-border rounded-full shadow-sm transition-all duration-300",
        "focus-within:border-neutral-300 focus-within:shadow-md", // Focus más sutil, sin aro amarillo feo
        className
      )}
    >
      <div className="flex items-center flex-1 w-full px-4 h-10">
        <span className="text-muted-foreground font-medium whitespace-nowrap mr-1 text-sm">
          uvico.shop/
        </span>
        <input
          type="text"
          value={slug}
          onChange={(e) => {
            const val = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "");
            setSlug(val);
          }}
          placeholder="mi-tienda"
          className="w-full bg-transparent outline-none placeholder:text-neutral-300 text-foreground font-medium text-sm"
        />
      </div>
      <Button 
        variant={variant} 
        className="w-full sm:w-auto mt-1.5 sm:mt-0 h-10 px-6 rounded-full text-sm"
      >
        Crear tienda
      </Button>
    </div>
  );
}

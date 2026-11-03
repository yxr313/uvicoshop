"use client";

import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { IconShieldCheck, IconStar, IconTruckDelivery } from "@tabler/icons-react";
import type { Supplier } from "./types";

export function SupplierCard({ supplier }: { supplier: Supplier }) {
  return (
    <Card className="flex flex-col h-full transition-all hover:border-foreground/20 hover:shadow-md">
      <CardContent className="p-5 flex flex-col h-full">
        
        <div className="flex items-start gap-4 mb-4">
          <div className="size-12 rounded-xl bg-neutral-100 flex items-center justify-center text-lg font-semibold text-neutral-600 shrink-0">
            {supplier.logoInitials}
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-1.5">
              <h3 className="font-semibold text-foreground truncate">{supplier.name}</h3>
              {supplier.isVerified && <IconShieldCheck className="size-4 text-[#95B54D] shrink-0" />}
            </div>
            <p className="text-xs text-muted-foreground mt-0.5">
              {supplier.category} • {supplier.state}
            </p>
          </div>
        </div>

        <p className="text-sm text-muted-foreground line-clamp-2 mb-4 flex-1">
          {supplier.description}
        </p>

        <div className="grid grid-cols-2 gap-2 mb-4 bg-neutral-50 p-2.5 rounded-lg border border-border/50">
          <div className="flex items-center gap-2 text-xs">
            <IconTruckDelivery className="size-4 text-muted-foreground" />
            <span className="text-foreground font-medium">{supplier.slaHours}h</span>
            <span className="text-muted-foreground">Despacho</span>
          </div>
          <div className="flex items-center gap-2 text-xs">
            <IconStar className="size-4 text-[#A7C957]" />
            <span className="text-foreground font-medium">{supplier.trustScore.toFixed(1)}</span>
            <span className="text-muted-foreground">Trust</span>
          </div>
        </div>

        <Button variant="outline" className="w-full border-border-strong" asChild>
          <Link href={`/dashboard/suppliers/dropshipping/${supplier.id}`}>
            Ver Perfil y Catálogo
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
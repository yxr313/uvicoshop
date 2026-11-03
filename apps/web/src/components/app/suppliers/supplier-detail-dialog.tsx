"use client";

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { IconShieldCheck, IconFileText, IconRefresh, IconShoppingBag, IconMapPin, IconReceipt } from "@tabler/icons-react";
import type { Supplier } from "./types";

interface SupplierDetailDialogProps {
  supplier: Supplier | null;
  isOpen: boolean;
  onClose: () => void;
  onConnect: (supplier: Supplier) => void;
}

export function SupplierDetailDialog({ supplier, isOpen, onClose, onConnect }: SupplierDetailDialogProps) {
  if (!supplier) return null;

  // Costo estimado de la muestra (ej. el producto más barato)
  const sampleCost = supplier.catalogPreview.length > 0 
    ? Math.min(...supplier.catalogPreview.map(p => p.priceCents)) 
    : 0;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-xl bg-white border-border shadow-xl">
        <DialogHeader>
          <div className="flex items-center gap-3 mb-2">
            <div className="size-14 rounded-xl bg-neutral-100 flex items-center justify-center text-xl font-semibold text-neutral-600">
              {supplier.logoInitials}
            </div>
            <div>
              <DialogTitle className="flex items-center gap-2 text-2xl">
                {supplier.name}
                {supplier.isVerified && <IconShieldCheck className="size-6 text-[#95B54D]" />}
              </DialogTitle>
              <DialogDescription className="flex items-center gap-2 mt-1">
                <span className="font-mono text-xs bg-neutral-100 px-2 py-0.5 rounded text-foreground">{supplier.rif}</span>
                <span>•</span>
                <span className="flex items-center gap-1"><IconMapPin className="size-3" /> {supplier.state}</span>
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-6 py-4">
          {/* Detalles Operativos */}
          <div className="grid grid-cols-2 gap-4">
            <div className="p-3 border border-border rounded-lg bg-neutral-50">
              <p className="text-xs text-muted-foreground mb-1">Dirección Física</p>
              <p className="text-sm font-medium text-foreground">{supplier.address}</p>
            </div>
            <div className="p-3 border border-border rounded-lg bg-neutral-50">
              <p className="text-xs text-muted-foreground mb-1">Tiempo de Despacho</p>
              <p className="text-sm font-medium text-foreground">{supplier.slaHours} horas (SLA)</p>
            </div>
          </div>

          {/* Catálogo de Muestra */}
          <div>
            <h4 className="text-sm font-medium text-foreground flex items-center gap-2 mb-3">
              <IconShoppingBag className="size-4 text-muted-foreground" /> Catálogo Destacado (Precios Base)
            </h4>
            <div className="grid grid-cols-2 gap-3">
              {supplier.catalogPreview.map((item) => (
                <div key={item.id} className="p-3 border border-border rounded-lg bg-white flex justify-between items-center">
                  <p className="text-sm font-medium text-foreground truncate">{item.name}</p>
                  <p className="text-sm text-muted-foreground font-mono">${(item.priceCents / 100).toFixed(2)}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Políticas Legales y Escrow */}
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 space-y-4">
            <div className="flex items-start gap-3">
              <IconReceipt className="size-5 text-amber-600 shrink-0" />
              <div>
                <p className="text-sm font-medium text-amber-900">Costo de Muestra Obligatorio</p>
                <p className="text-xs text-amber-800 mt-1">
                  Antes de activar la conexión, deberás comprar una muestra física. Costo estimado: <strong>${(sampleCost / 100).toFixed(2)}</strong>. Esto valida la calidad del proveedor.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <IconRefresh className="size-5 text-amber-600 shrink-0" />
              <div>
                <p className="text-sm font-medium text-amber-900">Garantía de Fondos (Escrow)</p>
                <p className="text-xs text-amber-800 mt-1">
                  Uvico retendrá los fondos de tus clientes hasta que el proveedor confirme el despacho y la guía de envío. Esto protege tu inversión.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <IconFileText className="size-5 text-amber-600 shrink-0" />
              <div>
                <p className="text-sm font-medium text-amber-900">Política de Devoluciones</p>
                <p className="text-xs text-amber-800 mt-1">{supplier.returnPolicy}</p>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter className="flex-col sm:flex-row gap-2">
          <Button variant="outline" onClick={onClose} className="w-full sm:w-auto">Cancelar</Button>
          <Button 
            className="bg-accent text-accent-foreground hover:bg-accent-hover w-full sm:w-auto"
            onClick={() => onConnect(supplier)}
          >
            Aceptar y Solicitar Muestra
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
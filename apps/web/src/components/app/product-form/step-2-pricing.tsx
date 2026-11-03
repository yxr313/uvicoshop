"use client";

import { useState } from "react";
import { useForm, Controller, type UseFormReturn } from "react-hook-form";
import type { ProductFormValues } from "@/hooks/use-product-form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { 
  IconCalculator, IconMinus, IconPlus, IconTruckDelivery, IconBox, 
  IconTrendingUp, IconCoin, IconPercentage, IconShoppingCart 
} from "@tabler/icons-react";
import { cn } from "@/lib/utils";

interface Step2Props {
  form: UseFormReturn<ProductFormValues>;
}

export function Step2Pricing({ form }: Step2Props) {
  const { control, watch, setValue, formState: { errors } } = form;
  const saleMode = watch("sale_mode");
  
  const price = watch("price_cents") || 0;
  const cost = watch("cost_cents") || 0;
  const shipping = watch("shipping_cost_cents") || 0;
  const ads = watch("ad_cost_cents") || 0;
  const stock = watch("initial_stock") || 0;

  // Estado para el simulador de dropshipping
  const [monthlySales, setMonthlySales] = useState(100);

  // --- LÓGICA FINANCIERA ---
  
  // 1. Flujo Inventario Propio (Cero comisiones)
  const totalOwnCost = cost + shipping + ads;
  const ownProfit = price - totalOwnCost;
  const ownMarginPercent = price > 0 ? (ownProfit / price) * 100 : 0;
  const ownRoi = totalOwnCost > 0 ? (ownProfit / totalOwnCost) * 100 : 0;
  const isOwnProfitable = ownProfit > 0;

  // 2. Flujo Dropshipping (Con comisiones Uvico y Pasarela)
  const uvicoFee = Math.round(price * 0.05);
  const gatewayFee = Math.round(price * 0.035);
  const dropshipProfit = price - cost - uvicoFee - gatewayFee;
  const isDropshipProfitable = dropshipProfit > 0;
  
  // Simulador
  const monthlyProfit = dropshipProfit * monthlySales;

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
      
      {/* SELECTOR DE MODO VISUAL RÁPIDO (Para no perder contexto) */}
      <div className="flex bg-neutral-100 p-1 rounded-xl w-fit">
        <button
          type="button"
          onClick={() => setValue("sale_mode", "own_inventory")}
          className={cn(
            "px-4 py-2 text-sm font-medium rounded-lg transition-all flex items-center gap-2",
            saleMode === "own_inventory" ? "bg-white shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
          )}
        >
          <IconBox className="size-4" /> Mi Inventario
        </button>
        <button
          type="button"
          onClick={() => setValue("sale_mode", "dropshipping_verified")}
          className={cn(
            "px-4 py-2 text-sm font-medium rounded-lg transition-all flex items-center gap-2",
            saleMode === "dropshipping_verified" ? "bg-white shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
          )}
        >
          <IconTruckDelivery className="size-4" /> Dropshipping / POD
        </button>
      </div>

      {/* FORMULARIO DE PRECIOS (Usando Controller para evitar el bug) */}
      <div className="grid gap-6 sm:grid-cols-2">
        <div className="space-y-2">
          <Label>Costo del producto <span className="text-muted-foreground font-normal">(Lo que pagas)</span></Label>
          <div className="relative">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground font-medium">$</span>
            <Controller
              control={control}
              name="cost_cents"
              render={({ field }) => (
                <Input
                  type="number"
                  step="0.01"
                  className="pl-8 h-12 text-lg font-medium bg-white"
                  value={field.value ? field.value / 100 : ""}
                  onChange={(e) => field.onChange(e.target.value ? Math.round(parseFloat(e.target.value) * 100) : 0)}
                  placeholder="0.00"
                />
              )}
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label>Precio de venta al público</Label>
          <div className="relative">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground font-medium">$</span>
            <Controller
              control={control}
              name="price_cents"
              render={({ field }) => (
                <Input
                  type="number"
                  step="0.01"
                  className="pl-8 h-12 text-lg font-medium bg-white border-foreground/30 focus-visible:ring-foreground"
                  value={field.value ? field.value / 100 : ""}
                  onChange={(e) => field.onChange(e.target.value ? Math.round(parseFloat(e.target.value) * 100) : 0)}
                  placeholder="0.00"
                />
              )}
            />
          </div>
          {errors.price_cents && <p className="text-xs text-destructive">{errors.price_cents.message}</p>}
        </div>
      </div>

      {/* ========================================================= */}
      {/* FLUJO 1: INVENTARIO PROPIO (Herramienta Financiera)       */}
      {/* ========================================================= */}
      {saleMode === "own_inventory" && (
        <div className="space-y-6 animate-in fade-in zoom-in-95">
          <div className="grid gap-6 sm:grid-cols-2 p-5 border border-border rounded-xl bg-white">
            <div className="space-y-2">
              <Label>Costo de Envío <span className="text-muted-foreground font-normal">(Por unidad)</span></Label>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground">$</span>
                <Controller
                  control={control}
                  name="shipping_cost_cents"
                  render={({ field }) => (
                    <Input
                      type="number"
                      step="0.01"
                      className="pl-8 bg-white"
                      value={field.value ? field.value / 100 : ""}
                      onChange={(e) => field.onChange(e.target.value ? Math.round(parseFloat(e.target.value) * 100) : 0)}
                      placeholder="0.00"
                    />
                  )}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Costo de Publicidad <span className="text-muted-foreground font-normal">(Por venta)</span></Label>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground">$</span>
                <Controller
                  control={control}
                  name="ad_cost_cents"
                  render={({ field }) => (
                    <Input
                      type="number"
                      step="0.01"
                      className="pl-8 bg-white"
                      value={field.value ? field.value / 100 : ""}
                      onChange={(e) => field.onChange(e.target.value ? Math.round(parseFloat(e.target.value) * 100) : 0)}
                      placeholder="0.00"
                    />
                  )}
                />
              </div>
            </div>
          </div>

          {/* Calculadora de Márgenes (Propio) */}
          <div className="bg-white rounded-xl border border-border shadow-sm overflow-hidden">
            <div className="grid grid-cols-1 md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-border">
              
              <div className="p-6 flex flex-col justify-center">
                <div className="flex items-center gap-2 text-muted-foreground mb-2">
                  <IconCoin className="size-4" /> <span className="text-sm font-medium">Inversión Total</span>
                </div>
                <p className="text-2xl font-semibold text-foreground">${(totalOwnCost / 100).toFixed(2)}</p>
                <p className="text-xs text-muted-foreground mt-1">Costo + Envío + Ads</p>
              </div>

              <div className="p-6 flex flex-col justify-center">
                <div className="flex items-center gap-2 text-muted-foreground mb-2">
                  <IconPercentage className="size-4" /> <span className="text-sm font-medium">Margen (ROI)</span>
                </div>
                <p className={cn("text-2xl font-semibold", isOwnProfitable ? "text-[#95B54D]" : "text-destructive")}>
                  {ownRoi.toFixed(1)}%
                </p>
                <p className="text-xs text-muted-foreground mt-1">Retorno sobre inversión</p>
              </div>

              <div className={cn(
                "p-6 flex flex-col justify-center",
                isOwnProfitable ? "bg-[#F0F5E6]/50" : "bg-red-50"
              )}>
                <div className="flex items-center gap-2 text-muted-foreground mb-2">
                  <IconTrendingUp className={cn("size-4", isOwnProfitable ? "text-[#95B54D]" : "text-destructive")} /> 
                  <span className="text-sm font-medium">Ganancia Libre</span>
                </div>
                <p className={cn("text-3xl font-display font-semibold", isOwnProfitable ? "text-[#95B54D]" : "text-destructive")}>
                  ${(ownProfit / 100).toFixed(2)}
                </p>
                <p className="text-xs text-muted-foreground mt-1">100% para ti (Sin comisiones)</p>
              </div>

            </div>
          </div>

          {/* Inventario */}
          <div className="space-y-3 p-5 border border-border rounded-xl bg-white flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <Label className="text-base">Inventario Físico</Label>
              <p className="text-sm text-muted-foreground">¿Cuántas unidades tienes listas para enviar?</p>
            </div>
            <div className="flex items-center gap-3">
              <Button type="button" variant="outline" size="icon" onClick={() => setValue("initial_stock", Math.max(0, stock - 1))}>
                <IconMinus className="size-4" />
              </Button>
              <Controller
                control={control}
                name="initial_stock"
                render={({ field }) => (
                  <Input
                    type="number"
                    className="text-center h-10 w-24 text-lg font-medium bg-white"
                    value={field.value || ""}
                    onChange={(e) => field.onChange(e.target.value ? parseInt(e.target.value) : 0)}
                  />
                )}
              />
              <Button type="button" variant="outline" size="icon" onClick={() => setValue("initial_stock", stock + 1)}>
                <IconPlus className="size-4" />
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* ========================================================= */}
      {/* FLUJO 2: DROPSHIPPING / POD (El Gancho Psicológico)       */}
      {/* ========================================================= */}
      {saleMode === "dropshipping_verified" && (
        <div className="space-y-6 animate-in fade-in zoom-in-95">
          
          <div className="space-y-3 p-5 border border-border rounded-xl bg-white">
            <div className="flex items-center gap-2">
              <IconTruckDelivery className="size-5 text-muted-foreground" />
              <Label className="text-base">Proveedor Conectado</Label>
            </div>
            <Select onValueChange={(val) => setValue("provider_id", val, { shouldValidate: true })}>
              <SelectTrigger className="h-12 bg-white">
                <SelectValue placeholder="Selecciona un proveedor" />
              </SelectTrigger>
              <SelectContent className="bg-white border-border">
                <SelectItem value="wayu_pay">Wayu Pay (Venezuela)</SelectItem>
                <SelectItem value="printful">Printful (POD Global)</SelectItem>
                <SelectItem value="cjdropshipping">CJ Dropshipping</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-xs text-muted-foreground">El stock y los envíos se sincronizan automáticamente.</p>
          </div>

          {/* Desglose de Dropshipping */}
          <div className="bg-white rounded-xl border border-border shadow-sm overflow-hidden">
            <div className="bg-neutral-50 px-6 py-4 border-b border-border flex items-center gap-3">
              <IconCalculator className="size-5 text-foreground" />
              <h3 className="font-medium text-foreground">Desglose de la Venta</h3>
            </div>
            
            <div className="p-6 grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="space-y-4">
                <div className="flex justify-between items-center text-sm">
                  <span className="text-muted-foreground">Precio de Venta</span>
                  <span className="font-medium text-foreground">${(price / 100).toFixed(2)}</span>
                </div>
                <div className="flex justify-between items-center text-sm">
                  <span className="text-muted-foreground">Costo Proveedor</span>
                  <span className="text-destructive">-${(cost / 100).toFixed(2)}</span>
                </div>
                <div className="flex justify-between items-center text-sm border-t border-dashed border-border pt-2">
                  <span className="text-muted-foreground">Comisión Uvico (5%)</span>
                  <span className="text-destructive">-${(uvicoFee / 100).toFixed(2)}</span>
                </div>
                <div className="flex justify-between items-center text-sm">
                  <span className="text-muted-foreground">Comisión Pasarela (3.5%)</span>
                  <span className="text-destructive">-${(gatewayFee / 100).toFixed(2)}</span>
                </div>
              </div>

              <div className={cn(
                "flex flex-col items-center justify-center rounded-lg border p-4",
                isDropshipProfitable ? "bg-[#F0F5E6] border-[#A7C957]/30" : "bg-red-50 border-red-200"
              )}>
                <p className="text-sm font-medium text-muted-foreground mb-1">Ganancia Neta por Venta</p>
                <p className={cn(
                  "text-4xl font-display font-semibold tracking-tight",
                  isDropshipProfitable ? "text-[#95B54D]" : "text-destructive"
                )}>
                  ${(dropshipProfit / 100).toFixed(2)}
                </p>
                {!isDropshipProfitable && price > 0 && (
                  <p className="text-xs text-destructive mt-2">Estás vendiendo a pérdida.</p>
                )}
              </div>
            </div>
          </div>

          {/* SIMULADOR DE ESCALA (El Gancho) */}
          {isDropshipProfitable && dropshipProfit > 0 && (
            <div className="bg-[#2736D9] text-white rounded-xl p-6 shadow-md flex flex-col md:flex-row items-center justify-between gap-6 animate-in slide-in-from-bottom-4">
              <div className="flex-1">
                <h3 className="text-xl font-medium mb-2 flex items-center gap-2">
                  <IconShoppingCart className="size-5" /> Simulador de Ganancias
                </h3>
                <p className="text-white/80 text-sm max-w-sm">
                  ¿Cuántas ventas estimas hacer al mes con este producto? Descubre tu potencial.
                </p>
              </div>
              
              <div className="flex flex-col items-center">
                <div className="flex items-center gap-2 bg-white/10 p-2 rounded-lg">
                  <input
                    type="number"
                    value={monthlySales}
                    onChange={(e) => setMonthlySales(Number(e.target.value) || 0)}
                    className="w-24 bg-transparent text-center text-2xl font-bold outline-none text-white border-b border-white/30"
                  />
                  <span className="text-sm font-medium text-white/80">ventas/mes</span>
                </div>
              </div>

              <div className="text-center md:text-right">
                <p className="text-sm text-white/80 mb-1">Ganancia Mensual Estimada</p>
                <p className="text-4xl font-display font-bold">
                  ${(monthlyProfit / 100).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
              </div>
            </div>
          )}

        </div>
      )}
    </div>
  );
}
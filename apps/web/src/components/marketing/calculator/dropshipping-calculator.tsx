"use client";

import { useState } from "react";
import { formatUsd, calculateSaleSplit } from "@uvico/pricing";
import { ArrowRight, AlertCircle } from "lucide-react";

export function DropshippingCalculator() {
  const [supplierCost, setSupplierCost] = useState<number>(4);
  const [salePrice, setSalePrice] = useState<number>(15);
  
  const gatewayRate = 0.035;
  const gatewayCost = salePrice * gatewayRate;
  const isValid = salePrice > supplierCost;

  let split = { uvicoCommissionUsd: 0, merchantNetUsd: 0 };
  if (isValid) {
    try {
      split = calculateSaleSplit({
        saleAmountUsd: salePrice,
        supplierCostUsd: supplierCost,
        gatewayFeeUsd: gatewayCost,
      });
    } catch (e) {}
  }

  return (
    <section className="py-32 px-4 bg-background">
      <div className="max-w-5xl mx-auto">
        
        <div className="text-center mb-20">
          <h2 className="text-display-2 text-foreground">
            Calcula tu margen real.
          </h2>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
          
          <div className="bg-neutral-50 p-8 md:p-10 rounded-[2rem] border border-border">
            <div className="mb-10">
              <div className="flex justify-between items-end mb-4">
                <span className="text-body font-medium text-foreground">Costo del proveedor</span>
                <span className="text-heading-lg text-foreground">{formatUsd(supplierCost)}</span>
              </div>
              <input 
                type="range" min="1" max="100" step="1"
                value={supplierCost}
                onChange={(e) => setSupplierCost(Number(e.target.value))}
                className="w-full h-2 bg-neutral-200 rounded-lg appearance-none cursor-pointer accent-neutral-400"
              />
            </div>

            <div className="mb-10">
              <div className="flex justify-between items-end mb-4">
                <span className="text-body font-medium text-foreground">Tu precio de venta</span>
                <span className="text-heading-lg text-foreground">{formatUsd(salePrice)}</span>
              </div>
              <input 
                type="range" min="5" max="200" step="1"
                value={salePrice}
                onChange={(e) => setSalePrice(Number(e.target.value))}
                className="w-full h-2 bg-neutral-200 rounded-lg appearance-none cursor-pointer accent-foreground"
              />
            </div>

            <div className="space-y-3 pt-6 border-t border-border">
              <div className="flex justify-between items-center text-body-sm">
                <span className="text-muted-foreground">Proveedor verificado</span>
                <span className="text-foreground font-medium">-{formatUsd(supplierCost)}</span>
              </div>
              <div className="flex justify-between items-center text-body-sm">
                <span className="text-muted-foreground">Comisión UVICO (5%)</span>
                <span className="text-foreground font-medium">-{formatUsd(isValid ? split.uvicoCommissionUsd : 0)}</span>
              </div>
              <div className="flex justify-between items-center text-body-sm">
                <span className="text-muted-foreground">Pasarela de pago (3.5%)</span>
                <span className="text-foreground font-medium">-{formatUsd(gatewayCost)}</span>
              </div>
            </div>
          </div>

          <div className="flex flex-col justify-center items-center lg:items-start text-center lg:text-left">
            <p className="text-body-lg text-muted-foreground mb-2">Tu ganancia neta</p>
            
            <div className="h-24 flex items-center mb-8">
              {isValid ? (
                <h3 className="text-[5rem] leading-none font-display font-medium tracking-tighter text-foreground">
                  {formatUsd(split.merchantNetUsd)}
                </h3>
              ) : (
                <div className="flex flex-col items-center lg:items-start">
                  <h3 className="text-[3rem] leading-none font-display font-medium tracking-tighter text-destructive mb-2">
                    Margen negativo
                  </h3>
                  <p className="text-body-sm text-destructive flex items-center gap-2">
                    <AlertCircle size={16} />
                    El precio de venta debe superar el costo.
                  </p>
                </div>
              )}
            </div>
            
            <button 
              className="bg-accent text-accent-foreground px-8 py-4 rounded-full font-medium transition-base hover:bg-accent-hover flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={!isValid}
            >
              Empezar a vender <ArrowRight size={18} />
            </button>
          </div>

        </div>
      </div>
    </section>
  );
}

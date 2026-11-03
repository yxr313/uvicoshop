"use client";

import { useState } from "react";
import type { UseFormReturn } from "react-hook-form";
import type { ProductFormValues } from "@/hooks/use-product-form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { IconPhoto, IconPlus, IconTrash } from "@tabler/icons-react";

interface Step3Props {
  form: UseFormReturn<ProductFormValues>;
}

export function Step3Gallery({ form }: Step3Props) {
  const { watch, setValue } = form;
  const hasVariations = watch("has_variations");
  const variations = watch("variations") || [];

  const [newVarName, setNewVarName] = useState("");
  const [newVarOptions, setNewVarOptions] = useState("");

  const addVariation = () => {
    if (!newVarName || !newVarOptions) return;
    const options = newVarOptions.split(",").map(s => s.trim()).filter(Boolean);
    setValue("variations", [...variations, { name: newVarName, options }], { shouldValidate: true });
    setNewVarName("");
    setNewVarOptions("");
  };

  const removeVariation = (index: number) => {
    const updated = [...variations];
    updated.splice(index, 1);
    setValue("variations", updated, { shouldValidate: true });
  };

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
      
      {/* Galería de Imágenes */}
      <div className="space-y-3">
        <Label className="text-base">Imágenes del producto</Label>
        <div className="border-2 border-dashed border-border-strong rounded-xl p-8 flex flex-col items-center justify-center text-center bg-neutral-50 hover:bg-neutral-100 transition-colors cursor-pointer">
          <div className="p-3 bg-white rounded-full shadow-sm mb-3">
            <IconPhoto className="size-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-foreground">Arrastra tus fotos aquí o haz clic para buscar</p>
          <p className="text-xs text-muted-foreground mt-1">JPG, PNG o WEBP. Máximo 5MB por foto.</p>
        </div>
      </div>

      {/* Variaciones */}
      <div className="space-y-4 border-t border-border pt-6">
        <div className="flex items-center justify-between">
          <div>
            <Label className="text-base">Variaciones</Label>
            <p className="text-sm text-muted-foreground">¿Este producto tiene diferentes tallas, colores o materiales?</p>
          </div>
          <Button 
            type="button" 
            variant={hasVariations ? "default" : "outline"}
            size="sm"
            onClick={() => setValue("has_variations", !hasVariations, { shouldValidate: true })}
          >
            {hasVariations ? "Quitar variaciones" : "Añadir variaciones"}
          </Button>
        </div>

        {hasVariations && (
          <div className="space-y-4 bg-neutral-50 p-4 rounded-xl border border-border">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3 items-end">
              <div className="space-y-1">
                <Label className="text-xs">Nombre (Ej. Talla)</Label>
                <Input value={newVarName} onChange={(e) => setNewVarName(e.target.value)} placeholder="Talla" />
              </div>
              <div className="space-y-1">
                <Label className="text-xs">Opciones separadas por coma (Ej. S, M, L)</Label>
                <Input value={newVarOptions} onChange={(e) => setNewVarOptions(e.target.value)} placeholder="S, M, L" />
              </div>
              <Button type="button" onClick={addVariation} variant="secondary">
                <IconPlus className="size-4 mr-2" /> Agregar
              </Button>
            </div>

            {variations.length > 0 && (
              <div className="border border-border rounded-lg bg-white overflow-hidden mt-4">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Variante</TableHead>
                      <TableHead>Opciones</TableHead>
                      <TableHead className="w-[100px]"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {variations.map((v, i) => (
                      <TableRow key={i}>
                        <TableCell className="font-medium">{v.name}</TableCell>
                        <TableCell>
                          <div className="flex gap-1 flex-wrap">
                            {v.options.map(opt => (
                              <span key={opt} className="bg-neutral-100 px-2 py-0.5 rounded text-xs">{opt}</span>
                            ))}
                          </div>
                        </TableCell>
                        <TableCell>
                          <Button variant="ghost" size="icon" onClick={() => removeVariation(i)}>
                            <IconTrash className="size-4 text-destructive" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
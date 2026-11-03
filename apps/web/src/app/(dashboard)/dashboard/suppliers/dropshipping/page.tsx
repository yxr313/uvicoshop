"use client";

import { useState } from "react";
import { PageHeader } from "@/components/page-header";
import { SupplierFilters } from "@/components/app/suppliers/supplier-filters";
import { SupplierCard } from "@/components/app/suppliers/supplier-card";
import type { Supplier } from "@/components/app/suppliers/types";

const MOCK_SUPPLIERS: Supplier[] = [
  {
    id: "textiles-maracay",
    name: "Textiles Maracay",
    description: "Fabricantes de ropa deportiva y casual de alta calidad. Envíos a todo el país.",
    category: "Ropa",
    state: "Carabobo",
    address: "Zona Industrial San Vicente, Galpón 4, Maracay",
    rif: "J-12345678-9",
    isVerified: true,
    trustScore: 4.8,
    slaHours: 24,
    logoInitials: "TM",
    connectionsCount: 145,
    productsCount: 24,
    completedOrders: 1250,
    responseTimeHours: 2,
    returnRate: 1.5,
    products: [
      { id: "camiseta-oversize", name: "Camiseta Oversize", description: "Algodón peinado 240gsm", priceCents: 800, rating: 4.9, monthlyRotation: 500, marketAdoption: 85 },
      { id: "franela-cuello", name: "Franela Cuello Redondo", description: "Corte clásico, ideal para estampados", priceCents: 650, rating: 4.7, monthlyRotation: 320, marketAdoption: 72 }
    ],
    documents: [
      { id: "d1", title: "Registro Mercantil", type: "Registro Mercantil", url: "#" },
      { id: "d2", title: "Políticas de Devolución", type: "Políticas", url: "#" }
    ]
  },
  {
    id: "calzados-caracas",
    name: "Calzados Caracas C.A.",
    description: "Zapatos formales y casuales para hombre y mujer. Excelente relación calidad-precio.",
    category: "Calzado",
    state: "Distrito Capital",
    address: "Av. Baralt, Esquina Puente Hierro, Caracas",
    rif: "J-98765432-1",
    isVerified: true,
    trustScore: 4.5,
    slaHours: 48,
    logoInitials: "CC",
    connectionsCount: 89,
    productsCount: 15,
    completedOrders: 840,
    responseTimeHours: 4,
    returnRate: 3.2,
    products: [
      { id: "zapato-formal", name: "Zapato Formal Cuero", description: "Cuero genuino, suela antideslizante", priceCents: 2500, rating: 4.6, monthlyRotation: 150, marketAdoption: 60 },
    ],
    documents: [
      { id: "d3", title: "Contrato de Suministro", type: "Contrato", url: "#" }
    ]
  }
];

export default function DropshippingDirectoryPage() {
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("all");
  const [state, setState] = useState("all");

  const filteredSuppliers = MOCK_SUPPLIERS.filter((s) => {
    const matchesSearch = s.name.toLowerCase().includes(search.toLowerCase()) || s.rif.includes(search);
    const matchesCategory = category === "all" || s.category === category;
    const matchesState = state === "all" || s.state === state;
    return matchesSearch && matchesCategory && matchesState;
  });

  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader title="Directorio de Proveedores" description="Explora y conecta con proveedores nacionales auditados por Uvico." />
      <SupplierFilters search={search} setSearch={setSearch} category={category} setCategory={setCategory} state={state} setState={setState} />

      {filteredSuppliers.length === 0 ? (
        <div className="p-12 flex flex-col items-center justify-center text-center min-h-[300px] border border-dashed rounded-xl bg-neutral-50">
          <h3 className="text-lg font-medium text-foreground mb-2">No se encontraron proveedores</h3>
          <p className="text-sm text-muted-foreground">Intenta ajustar los filtros o la búsqueda.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredSuppliers.map((supplier) => (
            <SupplierCard key={supplier.id} supplier={supplier} />
          ))}
        </div>
      )}
    </div>
  );
}
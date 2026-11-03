"use client";

import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { IconSearch, IconFilter } from "@tabler/icons-react";

interface SupplierFiltersProps {
  search: string;
  setSearch: (val: string) => void;
  category: string;
  setCategory: (val: string) => void;
  state: string;
  setState: (val: string) => void;
}

export function SupplierFilters({ search, setSearch, category, setCategory, state, setState }: SupplierFiltersProps) {
  return (
    <div className="flex flex-col lg:flex-row gap-3 w-full bg-white p-2.5 rounded-xl border border-border shadow-sm items-center">
      
      {/* Buscador */}
      <div className="relative flex-1 w-full">
        <IconSearch className="absolute left-3.5 top-1/2 -translate-y-1/2 text-muted-foreground size-4" />
        <Input
          placeholder="Buscar proveedor por nombre, RIF o ciudad..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-10 h-10 border-0 bg-transparent focus-visible:ring-0 shadow-none text-sm"
        />
      </div>

      {/* Filtros */}
      <div className="flex gap-2 items-center w-full lg:w-auto px-1 lg:border-l lg:border-border lg:pl-4">
        <IconFilter className="size-4 text-muted-foreground hidden sm:block shrink-0" />
        
        <Select value={category} onValueChange={setCategory}>
          <SelectTrigger className="w-full sm:w-auto min-w-[150px] h-10 bg-neutral-50 border-transparent text-sm">
            <SelectValue placeholder="Categoría" />
          </SelectTrigger>
          <SelectContent className="bg-white border-border">
            <SelectItem value="all">Todas las categorías</SelectItem>
            <SelectItem value="Ropa">Ropa</SelectItem>
            <SelectItem value="Calzado">Calzado</SelectItem>
            <SelectItem value="Tecnología">Tecnología</SelectItem>
            <SelectItem value="Hogar">Hogar</SelectItem>
            <SelectItem value="Belleza">Belleza</SelectItem>
            <SelectItem value="Deportes">Deportes</SelectItem>
          </SelectContent>
        </Select>

        <Select value={state} onValueChange={setState}>
          <SelectTrigger className="w-full sm:w-auto min-w-[170px] h-10 bg-neutral-50 border-transparent text-sm">
            <SelectValue placeholder="Estado" />
          </SelectTrigger>
          <SelectContent className="bg-white border-border">
            <SelectItem value="all">Toda Venezuela</SelectItem>
            <SelectItem value="Distrito Capital">Distrito Capital</SelectItem>
            <SelectItem value="Miranda">Miranda</SelectItem>
            <SelectItem value="Carabobo">Carabobo</SelectItem>
            <SelectItem value="Zulia">Zulia</SelectItem>
            <SelectItem value="Lara">Lara</SelectItem>
            <SelectItem value="Anzoátegui">Anzoátegui</SelectItem>
            <SelectItem value="Aragua">Aragua</SelectItem>
            <SelectItem value="Táchira">Táchira</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
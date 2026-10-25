"use client";
import { Puck } from "@puckeditor/core";
import "@puckeditor/core/puck.css";
import { puckConfig } from "@/lib/puck/puck.config";

export default function EditorPage() {
  const handlePublish = (data: any) => {
    console.log("Datos a guardar:", data);
    alert("Guardado localmente (revisa consola).");
  };
  return (
    <div className="h-[calc(100vh-8rem)] w-full rounded-xl overflow-hidden border border-border">
      <Puck config={puckConfig} data={{ content: [], root: {} }} onPublish={handlePublish} />
    </div>
  );
}

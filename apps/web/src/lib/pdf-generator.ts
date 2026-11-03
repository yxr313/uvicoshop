import { jsPDF } from "jspdf";
import autoTable from "jspdf-autotable";

interface PDFKpi {
  label: string;
  value: string;
  change: number;
}

interface PDFData {
  tenantName: string;
  range: string;
  kpis: PDFKpi[];
  paymentData: { source: string; value: number }[];
  topProducts: { name: string; quantity: number; revenue_cents: number }[];
  trafficData: { source: string; value: number }[];
}

// Extendemos la interfaz de jsPDF para que TypeScript reconozca `lastAutoTable`
interface ExtendedJsPDF extends jsPDF {
  lastAutoTable?: {
    finalY: number;
  };
}

export function generateAnalyticsPDF(data: PDFData) {
  const doc = new jsPDF() as ExtendedJsPDF;
  
  const primaryColor = [167, 201, 87]; // Verde Salvia Uvico
  const darkColor = [10, 10, 10];      // Negro Uvico
  const grayColor = [115, 115, 115];

  // 1. Marca de Agua
  doc.setTextColor(240, 240, 240);
  doc.setFontSize(60);
  doc.setFont("helvetica", "bold");
  doc.text("UVICO", 105, 150, { align: "center", angle: 45 });

  // 2. Encabezado
  doc.setTextColor(darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10);
  doc.setFontSize(24);
  doc.setFont("helvetica", "bold");
  doc.text("UVICO.", 14, 22);

  doc.setFontSize(16);
  doc.text("Inteligencia de Negocio", 14, 32);

  doc.setFontSize(10);
  doc.setFont("helvetica", "normal");
  doc.setTextColor(grayColor[0] ?? 115, grayColor[1] ?? 115, grayColor[2] ?? 115);
  doc.text(`Tienda: ${data.tenantName}`, 14, 40);
  doc.text(`Reporte generado el: ${new Date().toLocaleDateString("es-VE")}`, 14, 45);
  doc.text(`Período analizado: Últimos ${data.range}`, 14, 50);

  let yPos = 60;

  // 3. KPIs (Tabla)
  doc.setTextColor(darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10);
  doc.setFontSize(14);
  doc.setFont("helvetica", "bold");
  doc.text("Métricas de Comportamiento", 14, yPos);
  yPos += 5;

  autoTable(doc, {
    startY: yPos,
    head: [["Métrica", "Valor", "Tendencia"]],
    body: data.kpis.map((k) => [
      k.label,
      k.value,
      `${k.change > 0 ? "+" : ""}${k.change.toFixed(1)}%`,
    ]),
    theme: "grid",
    headStyles: { fillColor: [primaryColor[0] ?? 167, primaryColor[1] ?? 201, primaryColor[2] ?? 87], textColor: [darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10] },
    styles: { fontSize: 10, cellPadding: 4 },
  });

  yPos = (doc.lastAutoTable?.finalY ?? yPos) + 15;

  // 4. Métodos de Pago
  doc.setFontSize(14);
  doc.setFont("helvetica", "bold");
  doc.text("Métodos de Pago Preferidos", 14, yPos);
  yPos += 5;

  autoTable(doc, {
    startY: yPos,
    head: [["Método", "Transacciones"]],
    body: data.paymentData.map((p) => [p.source, p.value.toString()]),
    theme: "striped",
    headStyles: { fillColor: [darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10], textColor: [255, 255, 255] },
    styles: { fontSize: 10 },
  });

  yPos = (doc.lastAutoTable?.finalY ?? yPos) + 15;

  // 5. Top Productos
  doc.setFontSize(14);
  doc.setFont("helvetica", "bold");
  doc.text("Top Productos (Por Ingresos)", 14, yPos);
  yPos += 5;

  autoTable(doc, {
    startY: yPos,
    head: [["Producto", "Unidades", "Ingresos"]],
    body: data.topProducts.map((p) => [
      p.name,
      p.quantity.toString(),
      `$${(p.revenue_cents / 100).toLocaleString()}`,
    ]),
    theme: "grid",
    headStyles: { fillColor: [primaryColor[0] ?? 167, primaryColor[1] ?? 201, primaryColor[2] ?? 87], textColor: [darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10] },
    styles: { fontSize: 10 },
  });

  yPos = (doc.lastAutoTable?.finalY ?? yPos) + 15;

  // 6. Fuentes de Tráfico
  doc.setFontSize(14);
  doc.setFont("helvetica", "bold");
  doc.text("Fuentes de Tráfico", 14, yPos);
  yPos += 5;

  autoTable(doc, {
    startY: yPos,
    head: [["Canal", "Visitas / Registros"]],
    body: data.trafficData.map((t) => [t.source, t.value.toString()]),
    theme: "striped",
    headStyles: { fillColor: [darkColor[0] ?? 10, darkColor[1] ?? 10, darkColor[2] ?? 10], textColor: [255, 255, 255] },
    styles: { fontSize: 10 },
  });

  // 7. Footer (Número de página)
  const pageCount = doc.getNumberOfPages();
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i);
    doc.setFontSize(8);
    doc.setTextColor(grayColor[0] ?? 115, grayColor[1] ?? 115, grayColor[2] ?? 115);
    doc.text(
      `Página ${i} de ${pageCount} — Documento Confidencial de Uvico S.A.`,
      105,
      doc.internal.pageSize.height - 10,
      { align: "center" }
    );
  }

  // Guardar
  doc.save(`Reporte_UVICO_${new Date().toISOString().split("T")[0]}.pdf`);
}
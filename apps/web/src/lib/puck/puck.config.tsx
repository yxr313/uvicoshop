import type { Config } from "@puckeditor/core";

// Definimos los tipos de props para nuestros bloques
type Props = {
  HeadingBlock: { title: string };
  TextBlock: { text: string };
};

export const puckConfig: Config<Props> = {
  components: {
    HeadingBlock: {
      fields: {
        title: { type: "text" },
      },
      defaultProps: {
        title: "Bienvenido a Uvico",
      },
      render: ({ title }) => (
        <div style={{ padding: "32px 0" }}>
          <h1 style={{ fontSize: "2.5rem", fontWeight: "bold", margin: 0 }}>{title}</h1>
        </div>
      ),
    },
    TextBlock: {
      fields: {
        text: { type: "textarea" },
      },
      defaultProps: {
        text: "Este es el comienzo de tu tienda. Arrastra bloques para personalizarla.",
      },
      render: ({ text }) => (
        <div style={{ padding: "16px 0" }}>
          <p style={{ fontSize: "1.125rem", color: "#525252", lineHeight: "1.6" }}>{text}</p>
        </div>
      ),
    },
  },
};

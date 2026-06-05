import Atropos from "atropos/react";
import type { Card } from "../api/types";

type CardFaceProps = {
  card?: Card;
};

export function CardFace({ card }: CardFaceProps) {
  const hasFrontImage = Boolean(card?.image_url);

  return (
    <Atropos
      activeOffset={42}
      className={hasFrontImage ? "interactive-card interactive-card-selected" : "interactive-card"}
      duration={320}
      highlight
      rotateTouch="scroll-y"
      rotateXMax={10}
      rotateYMax={13}
      shadow={false}
      aria-label={card ? card.name : "Empty target card"}
    >
      <div className="interactive-card-flip">
        <div className="interactive-card-side interactive-card-back" aria-hidden={hasFrontImage}>
          <img className="interactive-card-artwork" src="/card-back.webp" alt="" data-atropos-offset="-1.5" />
        </div>
        <div className="interactive-card-side interactive-card-front">
          {card?.image_url ? <img className="interactive-card-artwork" src={card.image_url} alt={card.name} data-atropos-offset="-1.5" /> : null}
        </div>
      </div>
    </Atropos>
  );
}

export function CardFacts({ card }: { card: Card }) {
  const facts = buildCardFacts(card);

  return (
    <div className="fact-list" aria-label={`${card.name} facts`}>
      {facts.map((fact) => (
        <span key={fact}>{fact}</span>
      ))}
    </div>
  );
}

export function buildCardFacts(card: Card) {
  const stats = [];
  if (card.level !== null) {
    stats.push(`Level ${card.level}`);
  }
  if (card.rank !== null) {
    stats.push(`Rank ${card.rank}`);
  }
  if (card.link_rating !== null) {
    stats.push(`Link ${card.link_rating}`);
  }
  if (card.atk !== null || card.def !== null) {
    stats.push(`ATK ${card.atk ?? "?"} / DEF ${card.def ?? "?"}`);
  }

  const typeLine = buildTypeLine(card);
  const facts = [
    typeLine,
    card.attribute,
    ...stats,
    card.archetype ? `"${card.archetype}" card` : null,
    ...card.text_features.map((feature) => feature.split("_").join(" ")),
  ];

  return Array.from(new Set(facts.filter(Boolean) as string[]));
}

function buildTypeLine(card: Card) {
  if (card.card_type === "Monster") {
    const typeParts = [card.race, ...card.monster_categories].filter(Boolean);
    return typeParts.length ? `[${typeParts.join(" / ")}]` : "[Monster]";
  }

  if (card.card_type === "Spell" || card.card_type === "Trap") {
    const typeParts = [`${card.card_type} Card`, card.spell_trap_type].filter(Boolean);
    return `[${typeParts.join(" / ")}]`;
  }

  const typeParts = [card.card_type, card.frame_type].filter(Boolean);
  return typeParts.length ? `[${typeParts.join(" / ")}]` : null;
}

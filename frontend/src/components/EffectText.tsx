import type { SearcherResult } from "../api/types";

type HighlightKind = "condition" | "cost" | "action" | "restriction";

type Highlight = {
  kind: HighlightKind;
  text: string;
};

type Segment = {
  text: string;
  kind?: HighlightKind;
};

const highlightLabels: Record<HighlightKind, string> = {
  condition: "Condition",
  cost: "Cost",
  action: "Effect",
  restriction: "Restriction",
};

export function EffectText({
  hasTypeLine = false,
  result,
  fullText,
}: {
  hasTypeLine?: boolean;
  result: SearcherResult;
  fullText?: string;
}) {
  const segments = splitEffectText(result, fullText);

  return (
    <div className={hasTypeLine ? "effect-text" : "effect-text effect-text-standalone"}>
      <p>
        {segments.map((segment, index) =>
          segment.kind ? (
            <mark
              aria-label={`${highlightLabels[segment.kind]}: ${segment.text}`}
              className={`mark-${segment.kind}`}
              data-label={highlightLabels[segment.kind]}
              key={`${segment.kind}-${index}`}
              title={highlightLabels[segment.kind]}
            >
              {segment.text}
            </mark>
          ) : (
            <span key={`plain-${index}`}>{segment.text}</span>
          ),
        )}
      </p>
    </div>
  );
}

export function splitEffectText(result: SearcherResult, fullText?: string): Segment[] {
  const highlights: Highlight[] = [
    result.condition_text ? { kind: "condition", text: result.condition_text } : null,
    result.cost_text ? { kind: "cost", text: result.cost_text } : null,
    result.action_text ? { kind: "action", text: result.action_text } : null,
    result.restriction_text ? { kind: "restriction", text: result.restriction_text } : null,
  ].filter(Boolean) as Highlight[];

  const source = fullText ?? result.source_text;
  const matches = highlights
    .map((highlight) => {
      const index = source.indexOf(highlight.text);
      return index >= 0 ? { ...highlight, start: index, end: index + highlight.text.length } : null;
    })
    .filter(Boolean)
    .sort((left, right) => left!.start - right!.start || right!.end - left!.end) as Array<
    Highlight & { start: number; end: number }
  >;

  const segments: Segment[] = [];
  let cursor = 0;

  for (const match of matches) {
    if (match.start < cursor) {
      continue;
    }
    if (match.start > cursor) {
      segments.push({ text: source.slice(cursor, match.start) });
    }
    segments.push({ text: source.slice(match.start, match.end), kind: match.kind });
    cursor = match.end;
  }

  if (cursor < source.length) {
    segments.push({ text: source.slice(cursor) });
  }

  return segments.length > 0 ? segments : [{ text: source }];
}

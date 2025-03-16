export type TimeRange =
  | "Last month"
  | "Last week"
  | "Last 2 weeks"
  | "Last 3 days"
  | "Today"
  | "Yesterday";

export const dateAtUTC = (t: Date) => {
  //return new Date(t.getTime() + t.getTimezoneOffset() * 60000);
  return new Date(
    Date.UTC(
      t.getUTCFullYear(),
      t.getUTCMonth(),
      t.getUTCDate(),
      t.getUTCHours(),
      t.getUTCMinutes(),
      t.getUTCSeconds(),
    ),
  );
};

export const dateAtCurrentTZ = (t: string) => {
  const newTime = new Date(t);
  return newTime;
};

export const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(date);
};

export const getDateFromRange = (range: TimeRange): { from: Date; to: Date } => {
  const midnight = new Date();
  midnight.setHours(0, 0, 0, 0);
  const now = new Date();

  let t = { from: new Date(), to: new Date() };

  switch (range) {
    case "Last month":
      t = { from: new Date(midnight.getTime() - 30 * 24 * 60 * 60 * 1000), to: now };
      break;
    case "Last week":
      t = { from: new Date(midnight.getTime() - 7 * 24 * 60 * 60 * 1000), to: now };
      break;
    case "Last 2 weeks":
      t = { from: new Date(midnight.getTime() - 14 * 24 * 60 * 60 * 1000), to: now };
      break;
    case "Last 3 days":
      t = { from: new Date(midnight.getTime() - 3 * 24 * 60 * 60 * 1000), to: now };
      break;
    case "Yesterday":
      t = { from: new Date(midnight.getTime() - 24 * 60 * 60 * 1000), to: midnight };
      break;
    case "Today":
    default:
      t = { from: midnight, to: now };
  }

  return {
    from: dateAtUTC(t.from),
    to: dateAtUTC(t.to),
  };
};

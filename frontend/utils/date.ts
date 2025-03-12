export type TimeRange =
  | "Last month"
  | "Last week"
  | "Last 2 weeks"
  | "Last 3 days"
  | "Today"
  | "Yesterday";

export const dateAtUTC = (t: Date) => {
  return new Date(t.getTime() + t.getTimezoneOffset() * 60000);
};

export const dateAtCurrentTZ = (t: string) => {
  const newTime = new Date(t);
  return new Date(newTime.getTime() - newTime.getTimezoneOffset() * 60000);
};

export const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(date);
};

export const getDateFromRange = (range: TimeRange): { from: Date; to: Date } => {
  switch (range) {
    case "Last month":
      return {
        from: dateAtUTC(new Date(new Date().setMonth(new Date().getMonth() - 1))),
        to: dateAtUTC(new Date()),
      };
    case "Last week":
      return {
        from: dateAtUTC(new Date(new Date().setDate(new Date().getDate() - 7))),
        to: dateAtUTC(new Date()),
      };
    case "Last 2 weeks":
      return {
        from: dateAtUTC(new Date(new Date().setDate(new Date().getDate() - 14))),
        to: dateAtUTC(new Date()),
      };
    case "Last 3 days":
      return {
        from: dateAtUTC(new Date(new Date().setDate(new Date().getDate() - 3))),
        to: dateAtUTC(new Date()),
      };
    case "Today":
      return {
        from: dateAtUTC(new Date(new Date().setHours(0, 0, 0, 0))),
        to: dateAtUTC(new Date()),
      };
    case "Yesterday":
      return {
        from: dateAtUTC(new Date(new Date().setDate(new Date().getDate() - 1))),
        to: dateAtUTC(new Date(new Date().setHours(0, 0, 0, 0))),
      };
    default:
      return {
        from: dateAtUTC(new Date()),
        to: dateAtUTC(new Date()),
      };
  }
};

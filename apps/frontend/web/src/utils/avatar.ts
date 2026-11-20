export function getAvatarText(name: string) {
  const hasChinese = /[\u4e00-\u9fa5]/.test(name);

  if (hasChinese) {
    const chineseChars = name.match(/[\u4e00-\u9fa5]/g);
    return chineseChars
      ? chineseChars[chineseChars.length - 1]
      : name.slice(0, 2);
  }
  return name.replace(/\s/g, '').slice(0, 2).toUpperCase();
}

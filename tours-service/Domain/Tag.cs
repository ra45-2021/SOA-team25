using System.ComponentModel.DataAnnotations;

namespace TourService.Domain;

public class Tag
{
    public long Id { get; set; }

    [Required, MaxLength(64)]
    public string Name { get; set; } = string.Empty;

    public ICollection<TourTag> TourTags { get; set; } = new List<TourTag>();
}

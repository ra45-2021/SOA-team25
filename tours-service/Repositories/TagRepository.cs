using Microsoft.EntityFrameworkCore;
using TourService.Data;
using TourService.Domain;

namespace TourService.Repositories;

public class TagRepository : ITagRepository
{
    private readonly ToursContext _db;
    public TagRepository(ToursContext db) { _db = db; }

    public async Task<List<Tag>> GetOrCreateAsync(IEnumerable<string> tagNames, CancellationToken ct)
    {
        var normalized = tagNames
            .Where(x => !string.IsNullOrWhiteSpace(x))
            .Select(x => x.Trim())
            .Where(x => x.Length > 0)
            .Distinct(StringComparer.OrdinalIgnoreCase)
            .ToList();

        if (normalized.Count == 0) return new List<Tag>();

        var existing = await _db.Tags
            .Where(t => normalized.Contains(t.Name))
            .ToListAsync(ct);

        var existingSet = new HashSet<string>(existing.Select(x => x.Name), StringComparer.OrdinalIgnoreCase);

        var toAdd = normalized
            .Where(n => !existingSet.Contains(n))
            .Select(n => new Tag { Name = n })
            .ToList();

        if (toAdd.Count > 0)
        {
            _db.Tags.AddRange(toAdd);
            await _db.SaveChangesAsync(ct);
            existing.AddRange(toAdd);
        }

        return existing;
    }
}

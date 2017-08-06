class Plugin
  attr_accessor :context

  def self.generate(name)
    p = config(name) do |c|
      yield(c)
    end

    p.compile
  end

  def self.config(name)
    plugin = new(name)
    plugin.context = PluginContext.new
    yield(plugin.context)
    plugin
  end

  class PluginContext
    attr_accessor :functions, :packages

    def initialize
      @functions = []
      @packages = []
    end

    def link_function(prefix, name)
      @functions.push({ prefix: prefix, name: name })
    end

    def import_pkg(prefix, name)
      @packages.push({ prefix: prefix, name: name })
    end
  end
end